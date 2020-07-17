package endpoints

import "unsafe"

type Student struct {
	Name  string
	Age   int
	Class string
	Score int
}

func DirectInvoke(s *Student) {
	s.Name = "Jerry"
	s.Age = 18
	s.Class = "20005"
	s.Score = 100
}

func PointerInvoke(p unsafe.Pointer)  {
	s := (*Student)(p)
	s.Name = "Jerry"
	s.Age = 18
	s.Class = "20005"
	s.Score = 100
}
func InterfaceInvoke(i interface{}) {
	s := i.(*Student)
	s.Name = "Jerry"
	s.Age = 18
	s.Class = "20005"
	s.Score = 100
}
func InterfaceInvoke1(i interface{}) {
	s := i.(Student)
	s.Name = "Jerry"
	s.Age = 18
	s.Class = "20005"
	s.Score = 100
}