package time

/*
Defines common methods for a logical timestamp
NOTE: 	Due to the limitations of GoLang then a sync method is not defined in this interface. 
		When generics becomes available then the method should be defined.
*/
type LogicalTimestamp interface {
	Increment()
	GetDisplayableContent() string
}