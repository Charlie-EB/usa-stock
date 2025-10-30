
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

# packages
a package is a group of files in the same folder. and when the compiler does its thing, at the end there is one big file with all packages inside. seperating into files is just for our organisation. "code sees packages, not files"
packages can be: 
    build in standard lib package
    our own package
    a 3rd party dependency lib that we bring in

    and dont forget that package main in mandatory

dont need to import a user created package- its within the same folder so no need to import. if i want to have a "module" that i import around, then the user created func needs to be in a different folder (different package)

so question: in this root dir what package does the test.go belong to?

so every function etc. is being added to the main package. we're just extending it? not quite
when i write package main- im saying that this package will build an executable program, and go is expecting the func main entry point

interesting side note: print() is not guarateed to work in all oses
so we can use fmt.Print()
and wheter or not a function is private or not the naming convention is Print with an uppercase P for public. (we dont have the public keyword in go, just naming convention). Lowercase is private to the package and is not exported (made avaliable for import)




