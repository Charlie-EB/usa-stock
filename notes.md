
constants can only be bool, string or numbers.

why?
immutable variable != constant.

immutable var = address in memory that has to be alocated.
once set it cant be changed

and in JS at runtime this const is checked if its the same value

go is complied? ".. when the compiler finds the const reference, it goes to the value and copy paste it in. .. a fixed value"
There is an exception when it comes to strings
"..they wont change after compliation. const value set at compile time, not at runtime"

what does compile mean with respect to go?



# multiplatform
go is compiled and it needs to target a os and architecture (intel and arm)


what is := ? its a shorthand used to delcar a new variable. same as var newVar String = "hello world"
newVar := "hello world"
there is no shorthand for const. const 

data types:
uints are postive only
and can work with pointers. in a nutshell. variable and another variable pointing to a another's value
"a pointer is a variable that references the memory address of another variable". Link to original variable and allows for indirect access and manipulation of the orignal variable 
and every number in JSON is convereted to float64