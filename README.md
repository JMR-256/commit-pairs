# Commit Pairs 
A tool to wrap git functionality allowing for easier addition of co-authors to your commit messages

## Setup:

1. Whilst in the project root compile with ```go build .``` this will create the binary git-pc
2. Move the executable binary to a location on your PATH e.g. ```mv ./git-pc ~/bin```
3. Setup your .pairs file 
   - In your home directory create a file called .pairs e.g. ```touch ~/.pairs```
   - Note that currently only one domain can be supported at a time

Sample pairs file:

```
pairs:
  jr: James Riley; james.riley
  jd: John Doe; john.doe
  as: Alice Smith; alice.smith
  tj: Tom Jones; tom.jones
  
email:
  domain: example.com
```

#### You are now ready to use commit pairs!

## Usage:

Assuming the above setup process has been followed
```
git pc [primary initials] [co author initials]
```

Will allow you to quickly write commits as a user with any number of co-authors from your organisation

Example:
```
git pc jr tj as jd
```
With a valid .pairs file this will do the following: \
Set git config global user.name to James Riley \
Set git config global user.email to james.riley@example.com \
Create a .commitPairsTemplate file to be used for writing commits message with native text editor in the following format:
```

Co-authored-by: John Doe <john.doe@example.com>
Co-authored-by: Alice Smith <alice.smith@example.com>
Co-authored-by: Tom Jones <tom.jones@example.com>
```

Optionally run as script with Go:

```go run . jr ```