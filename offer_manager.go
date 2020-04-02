package main

import (
	"errors"
	"ets2-sync/db"
	"ets2-sync/dlc"
	"ets2-sync/savefile"
	"ets2-sync/structs"
	"ets2-sync/utils"
	"fmt"
	"math/rand"
	"sync"
	"time"
	"xorm.io/builder"
)

var currentOffers map[string][]structs.ApplicableOffer // key is SourceCompany
var totalOffersForSync int
var lastUpdatedSync time.Time

var offersInDb = make([]string, 0)
var jobToProcessMutex sync.Mutex
var createNewJobsMutex sync.Mutex
var jobsToProcess = make([]db.DbOffer, 0)
var jobManagerInitialized bool


func initOfferManager() error {
	if jobManagerInitialized {
		return errors.New("job manager is already initialized")
	}

	go func() {
		for {
			jobToProcessMutex.Lock()
			jobsCopy := jobsToProcess
			jobsToProcess = make([]db.DbOffer, 0)
			jobToProcessMutex.Unlock()

			batchSize := 100
			lenJobs := len(jobsCopy)

			context := db.GetDb()

			for i := 0; i < lenJobs; i += batchSize {
				j := i + batchSize

				if j > lenJobs {
					j = lenJobs
				}

				currentOffers := jobsCopy[i:j]
				hashes := make([]string, 0)

				for _, offer := range currentOffers {
					hashes = append(hashes, offer.Id)
				}

				createNewJobsMutex.Lock()
				realOffersInDb := make([]db.DbOffer, 0)

				_ = context.Where(builder.In("hash", hashes)).
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
	context := db.GetDb()

	currentOffersArr := make([]*db.DbOffer, 0)

	if err := context.OrderBy("required_dlc asc").Find(&currentOffersArr); err != nil {
		return err
	}

	currentOffers = make(map[string][]structs.ApplicableOffer)
	tempOffers := make(map[string][]*db.DbOffer)
	totalOffersForSync = 0
	lastUpdatedSync = time.Now().UTC()
	maxNonDlcJobs := 5
	maxCargoThreshold := 20
	segment3 := 0
	segment4 := 0

	for _, offer := range currentOffersArr {
		val, ok := tempOffers[offer.SourceCompany]
		offersInDb = append(offersInDb, offer.Id)

		if !ok {
			val = make([]*db.DbOffer, 0)
		}

		tempOffers[offer.SourceCompany] = append(val, offer)
	}

	for key, offers := range tempOffers {

		rand.Shuffle(len(offers), func(i, j int) {
			offers[i], offers[j] = offers[j], offers[i]
		})

		finalOffers := make([]*db.DbOffer, 0)

		for _, offer := range offers {
			if maxNonDlcJobs <= len(finalOffers) {
				break
			}

			if offer.RequiredDlc == dlc.BaseGame {
				finalOffers = append(finalOffers, offer)
			}
		}

		for _, offer := range offers {
			if offer.RequiredDlc == dlc.BaseGame {
				continue
			}

			if maxCargoThreshold <= len(finalOffers) {
				break
			}

			finalOffers = append(finalOffers, offer)
		}

		result := make([]structs.ApplicableOffer, 0)

		for _, offer := range finalOffers {
			result = append(result, structs.NewApplicableOffer(offer, fmt.Sprintf("_nameless.19a.%04d.%04d", segment3, segment4)))
			segment4++

			if segment4 == 9999 {
				segment4 = 0
				segment3++
			}
		}

		totalOffersForSync += len(finalOffers)

		currentOffers[key] = result
	}

	return nil
}

func PopulateOffers(file *savefile.SaveFile, supportedDlc dlc.Dlc) {
	for _, offer := range getOffers(supportedDlc, file.AvailableCompanies) {
		if err := file.AddOffer(offer); err != nil {
			fmt.Println(err) // todo
		}
	}
}

func getOffers(supportedDlc dlc.Dlc, availableSources []string) []structs.ApplicableOffer {
	offersToAdd := make([]structs.ApplicableOffer, 0)
	for sourceCompany, offers := range currentOffers {
		if !utils.Contains(availableSources, sourceCompany) {
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

func FillDbWithJobs(offers []*savefile.JobOffer) {
	if offers == nil || len(offers) == 0 {
		return
	}

	go func() {
		for _, v := range offers {
			offer := v.ToDbOffer()

			offer.RequiredDlc = dlc.GetRequiredDlc(v.SourceCompany, v.Target,
				v.Cargo, v.TrailerDefinition, v.TrailerVariant)

			if offer.RequiredDlc == dlc.None {
				continue // skip unverified offers
			}

			offer.Id = offer.CalculateHash()

			if len(offer.Id) == 0 {
				continue
			}

			if utils.Contains(offersInDb, offer.Id) {
				continue
			}

			jobToProcessMutex.Lock()

			shouldAdd := true

			for _, j := range jobsToProcess {
				if j.Id == offer.Id {
					shouldAdd = false
					break
				}
			}

			if !shouldAdd {
				jobToProcessMutex.Unlock()
				continue
			}

			jobsToProcess = append(jobsToProcess, offer)
			jobToProcessMutex.Unlock()
		}
	}()
}
