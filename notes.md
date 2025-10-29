
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