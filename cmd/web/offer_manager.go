package web

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"ets2-sync/dlc_mapper"
	"ets2-sync/internal"
	"ets2-sync/savefile"
	"xorm.io/builder"
)

var currentOffers = map[dlc_mapper.Game]map[string][]savefile.ApplicableOffer{dlc_mapper.ETS: nil, dlc_mapper.ATS: nil} // key is SourceCompany
var totalOffersForSync = map[dlc_mapper.Game]int{dlc_mapper.ETS: 0, dlc_mapper.ATS: 0}
var lastUpdatedSync time.Time

var offersInDb = make([]string, 0)
var jobToProcessMutex sync.Mutex
var createNewJobsMutex sync.Mutex
var jobsToProcess = map[dlc_mapper.Game][]savefile.ApplicableOffer{
	dlc_mapper.ETS: make([]savefile.ApplicableOffer, 0),
	dlc_mapper.ATS: make([]savefile.ApplicableOffer, 0),
}
var jobManagerInitialized bool

func initOfferManager() error {
	if jobManagerInitialized {
		return errors.New("job manager is already initialized")
	}

	go func() {
		for {
			for _, game := range dlc_mapper.AllGames {
				jobToProcessMutex.Lock()

				dbOffers := make([]dbOffer, 0)
				var ids []string
				for _, offer := range jobsToProcess[game] {
					if internal.Contains(ids, offer.Id) {
						continue
					}
					ids = append(ids, offer.Id)
					dbOffer := dbOffer{}
					_, _ = internal.MapToObject(offer, &dbOffer)

					dbOffers = append(dbOffers, dbOffer)
				}

				jobsToProcess[game] = make([]savefile.ApplicableOffer, 0)
				jobToProcessMutex.Unlock()

				batchSize := 100
				lenJobs := len(dbOffers)

				context := GetDb()

				for i := 0; i < lenJobs; i += batchSize {
					j := i + batchSize

					if j > lenJobs {
						j = lenJobs
					}

					currentOffers := dbOffers[i:j]
					ids := make([]string, 0)

					for _, offer := range currentOffers {
						ids = append(ids, offer.Id)
					}

					createNewJobsMutex.Lock()
					realOffersInDb := make([]dbOffer, 0)

					_ = context.Where(builder.In("id", ids)).
						Find(&realOffersInDb)

					for _, dbOffer := range realOffersInDb { // ideally it should never happen
						for _, offer := range currentOffers {
							if offer.Id == dbOffer.Id {
								offer.Id = "" // we need to skip that shit
							}
						}
					}

					for _, offer := range currentOffers {
						if len(offer.Id) == 0 {
							continue
						}

						if _, er := context.Insert(offer); er != nil {
							fmt.Println("shit") // todo log error
						} else {
							offersInDb = append(offersInDb, offer.Id)
						}
					}

					createNewJobsMutex.Unlock()
				}
			}

			time.Sleep(3 * time.Second) // todo config
		}
	}()

	go func() {
		for {
			_ = updateList()
			time.Sleep(time.Hour * 72)
		}
	}()

	return nil
}

func updateList() error {
	context := GetDb()
	currentOffers[dlc_mapper.ETS] = make(map[string][]savefile.ApplicableOffer)
	currentOffers[dlc_mapper.ATS] = make(map[string][]savefile.ApplicableOffer)
	currentOffers[dlc_mapper.ETS_PROMODS] = make(map[string][]savefile.ApplicableOffer)

	for _, game := range dlc_mapper.AllGames {

		currentOffersArr := make([]*dbOffer, 0)

		if err := context.Where("game = ?", game).OrderBy("required_dlc asc").Find(&currentOffersArr); err != nil {
			return err
		}

		tempOffers := make(map[string][]*dbOffer)
		totalOffersForSync[game] = 0
		lastUpdatedSync = time.Now().UTC()
		maxNonDlcJobs := 5
		maxCargoThreshold := 20
		segment3 := 0
		segment4 := 0

		for _, offer := range currentOffersArr {
			val, ok := tempOffers[offer.SourceCompany]
			offersInDb = append(offersInDb, offer.Id)

			if !ok {
				val = make([]*dbOffer, 0)
			}

			tempOffers[offer.SourceCompany] = append(val, offer)
		}

		for key, offers := range tempOffers {

			rand.Shuffle(len(offers), func(i, j int) {
				offers[i], offers[j] = offers[j], offers[i]
			})

			finalOffers := make([]*dbOffer, 0)

			for _, offer := range offers {
				if maxNonDlcJobs <= len(finalOffers) {
					break
				}

				if offer.RequiredDlc == dlc_mapper.BaseGame {
					finalOffers = append(finalOffers, offer)
				}
			}

			for _, offer := range offers {
				if offer.RequiredDlc == dlc_mapper.BaseGame {
					continue
				}

				if maxCargoThreshold <= len(finalOffers) {
					break
				}

				finalOffers = append(finalOffers, offer)
			}

			result := make([]savefile.ApplicableOffer, 0)

			for _, offer := range finalOffers {
				newOffer := savefile.ApplicableOffer{}
				_, _ = internal.MapToObject(offer, &newOffer)

				newOffer.Id = fmt.Sprintf("_nameless.19a.%04d.%04d", segment3, segment4)

				result = append(result, newOffer)
				segment4++

				if segment4 == 9999 {
					segment4 = 0
					segment3++
				}
			}

			totalOffersForSync[game] += len(finalOffers)

			currentOffers[game][key] = result
		}
	}

	return nil
}

func PopulateOffers(file *savefile.SaveFile, game dlc_mapper.Game, supportedDlc dlc_mapper.Dlc) {
	for _, offer := range getOffers(game, supportedDlc, file.AvailableCompanies) {
		if err := file.AddOffer(offer); err != nil {
			fmt.Println(err) // todo
		}
	}
}

func getOffers(game dlc_mapper.Game, supportedDlc dlc_mapper.Dlc, availableSources []string) []savefile.ApplicableOffer {
	offersToAdd := make([]savefile.ApplicableOffer, 0)
	for sourceCompany, offers := range currentOffers[game] {
		if !internal.Contains(availableSources, sourceCompany) {
			continue
		}

		for _, offer := range offers {
			if (supportedDlc & offer.RequiredDlc) == offer.RequiredDlc {
				offersToAdd = append(offersToAdd, offer)
			}
		}
	}

	return offersToAdd
}

func FillDbWithJobs(offers []savefile.ApplicableOffer, game dlc_mapper.Game) {
	if offers == nil || len(offers) == 0 {
		return
	}

	go func() {
		for _, offer := range offers {
			offer.RequiredDlc = dlc_mapper.GetRequiredDlc(offer.SourceCompany, offer.Target,
				offer.Cargo, offer.TrailerDefinition, offer.TrailerVariant, game)
			offer.Game = game

			if offer.RequiredDlc == dlc_mapper.None {
				continue // skip unverified offers
			}

			offer.Id = offer.CalculateHash()

			if len(offer.Id) == 0 {
				continue
			}

			if internal.Contains(offersInDb, offer.Id) {
				continue
			}

			jobToProcessMutex.Lock()

			shouldAdd := true

			for _, j := range jobsToProcess[game] {
				if j.Id == offer.Id {
					shouldAdd = false
					break
				}
			}

			if !shouldAdd {
				jobToProcessMutex.Unlock()
				continue
			}

			jobsToProcess[game] = append(jobsToProcess[game], offer)
			jobToProcessMutex.Unlock()
		}
	}()
}
