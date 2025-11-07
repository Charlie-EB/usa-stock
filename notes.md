# notes about learning go
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


collection
arrays in go have a fixed length. denoted by [5]int. array of length 5, type int
slices are similar to array with a dynamic length. denoted by []int. 
actually behind the scenes we're looking at a slice/view of a real fixed length array. slices have a dynamic length. when adding elements to a slice, a new slice is created
maps key value. denoted map[int]string

reminder: collections are not objects. nothing is an object. 
so we use global functions like len() on the collection to get the length. 
not array.length()



side note:
i can delcare many init() functions with the same name


functions
can return more than 1 value. not an collection, but 2 values

pointers
funcs by default have arguments passed by value ie a copy of the original value.
we can also pass in a pointer to the original value / reference to the value and modify it.
e.g.
func birthday(pointerAge *int) {
    *pointerAge++
    }

func main() {
    age:= 22
    birthday(&age)
    print(age)
    print(&age)
    }

* is the pointer
& is passing in the reference (memory address to the pointer)

and in the above example, we're passing in a memory address. and telling the function to accept the address and increment it via * 

# errors
dont have an error func, so there is design pattern

func readUser(id int) (user, err) {
    // we're reading something and see a bool ok 
    if ok {
        return user, nil
    } else {
        return nil, errorDetails
        }
func main {
    user, err := readUser(2)
    }
            

# control flow
go only has one equality operator ==



    if else - can delcare a top level variable that is avaliable in each local if and else block
    switch - the case can be evauled to a bool can replace large if else blocks, and we have a fallthrough value. breaks by default
    and switches execult only on the first true case

    for - can use classic for loop. set i; i < length; i++
    or can use for in
    for index := range collection {}
    and can use for each
    for key, value := range map {}

    and can replicate while. using a statement that evalutes to a bool
    in summary for loops can be used via classic, range iterator over collection, and directly witha  boolena expression (replacing while loop)

