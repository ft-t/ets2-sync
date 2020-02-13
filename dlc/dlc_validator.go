package dlc

type IDLCValidator interface {
	CanAdd(string) bool
}

type BaseGameValidator struct {

}
