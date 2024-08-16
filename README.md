# Commit Pairs

A tool to wrap Git functionality, allowing for easier addition of co-authors to your commit messages.

## Setup

1. **Compile the Project:**
   - While in the project root, compile the application using:
     ```sh
     go build .
     ```
     This will create the binary `git-pc`.

2. **Move the Executable to Your PATH:**
   - Move the executable binary to a location on your `PATH`, for example:
     ```sh
     mv ./git-pc ~/bin
     ```

3. **Set Up Your `.pairs` File:**
   - In your home directory, create a file named `.pairs`:
     ```sh
     touch ~/.pairs
     ```
   - Note that currently, only one domain can be supported at a time.

### Sample `.pairs` File:

```yaml
pairs:
  jr: James Riley; james.riley
  jd: John Doe; john.doe
  as: Alice Smith; alice.smith
  tj: Tom Jones; tom.jones
  
email:
  domain: example.com
```

#### You are now ready to use Commit Pairs!

## Usage

After completing the setup process, you can use the following command:

```sh
git pc [primary initials] [co-author initials...]
```

This allows you to quickly create commits as a user with any number of co-authors from your organization.

### Example:

```sh
git pc jr tj as jd
```

With a valid `.pairs` file, this command will:

1. Set the global Git config `user.name` to **James Riley**.
2. Set the global Git config `user.email` to **james.riley@example.com**.
3. Create a `.commitPairsTemplate` file for writing commit messages with your native text editor. The commit message template will include the following format:

    ```text
    Co-authored-by: John Doe <john.doe@example.com>
    Co-authored-by: Alice Smith <alice.smith@example.com>
    Co-authored-by: Tom Jones <tom.jones@example.com>
    ```

### Running as a Script with Go:

If you prefer to run the tool without building, you can execute it directly with Go:

```sh
go run . jr
```