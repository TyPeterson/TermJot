# TermJot
---

## Overview
TermJot is a simple command line tool (CLI) for quickly jotting down notes or 'terms' from within the terminal. The implementation is written in Go, and uses the Cobra CLI library for handling the CLI commands and flags. The notes are stored in a SQLite database, and the database file is created and stored in the user's home directory as a hidden file (.termjot.db).

* Notes/terms are organized by 'Category', which is defined by the user at the time of note creation (if no category is given, the note is added to the global 'ALL' category).

* In addition to giving a note a category, the user can also optionally add a 'Definition' to a note, if the note is meant to be a term that the user wishes to remember such as a specific keyboard shortcut, or possibly the syntax to do some specific thing in a programming language. An example use case of this could be having notes within the 'Git' category that contain shorthand notes on how to do specific things with git.

* TermJot also allows you to quickly ask any question and have an answered generated by the gemini-1.5-pro model from Google. This is done by using the 'ask' command, followed by the question (and option -c flag, followed by the context/category of the question).

Because of the emphasis on quick and easy access from within your existing dev environment, some features were added to enhance speed and convenience:
  1. When providing a category, you can instead just type '.' to use the current project directory as the category name. This is useful for keeping track of any TODO items or notes for a project that you are currently working on. For example, while working on this tool, I was able to quicky add a TODO item by saying, `  tjot add .  ` and then giving the content of the note, which would be added to the 'TermJot' category/list.

  2. When asking a question with the 'ask' command, you can optionally provide the -b (brief) flag to have the answer given in as few words as possible. For example, if I just wanted to quickly remember the vim shortcut for searching and replacing text throughout a file, I could say, `  tjot ask -b 'vim search and replace'  ` and in < 1 second I would get the answer, `  :%s/search/replace/g  `.

  3. Quickly add a file contents to the prompt with the -f flag in the 'ask' command. This will pass the entire contents of the given file to the prompt to use for context.

  4. Default model is gemini-1.5-flash since it is the best current balance of insane speed, near SOTA performance, and as a bonus, free to use. Since the -b flag will result in the reponse using as few words as possible, the model only ever has to generate a couple dozen tokens before reaching the <END_OF_RESPONSE> token, which is when the reponse is then sent back to the user. This means that the gemini-1.5-flash model, which generates ~140 tokens/s, will lead to a response being shown to the user in < 1 second when the -b flag is used. Standard responses will take ~2 seconds, which is still very fast.
---

## Installation
To use this tool you will need to have Go installed on your machine. If you do not have Go installed, you can download it from the official Go website [here](https://golang.org/dl/).

Sadly, on macOS, 'jot' is a built-in system utility, so you cannot name the output binary 'jot'. This is why I name my output binary 'tjot', but if you use Linux or Windows, you can freely name the output binary 'jot' for a shorter name to have to type to use the CLI tool.

Once you have Go installed, you can clone this repository with:
```
git clone http://github.com/TyPeterson/TermJot.git
```

Then you can navigate into the TermJot directory and build the binary with:
```
cd TermJot
go build -o tjot
```

Finally, you can move the binary to your $PATH with:

> **macOS & Linux:**
> ```
> sudo mv tjot /usr/local/bin
> ```
> and possibly may need to ensure the binary has the correct permissions with:
> ```
> sudo chmod +x /usr/local/bin/tjot
> ```

or

> **Windows:**
> 1. copy the binary to C:\Program Files\TermJot with:
> ```
> copy tjot.exe "C:\Program Files\TermJot"
> ```
> 2. Edit the system environment variables and edit the Path variable to include "C:\Program Files\TermJot"

verify that the binary is not globally accessible with:
```
tjot
```
---
## Usage
TermJot supports the following commands:

### Add
> *Add a new term to the global 'ALL' category, or optionally to a specific category.*

> **Arguments:**
> * category: (optional - string) The category that the term/definition will be added within. If no category is given, the default category is 'ALL'.

> **Flags:**
> * -t, --term: (optional - string) The term to add to the database, surrounded by quotes.
> * -d, --define: (optional - bool) Use this flag if you want to add a definition to an existing term. If given, the user will then select the term to add the definition to, followed by providing the definition.

*Note: the -t and -d flags are mutually exclusive, and only one can be used at a time.*


### List
> *List all terms within all categories, or if a category is given, list all terms within that category.*<br>

> **Arguments:**
> * category: (optional - string) The category to list terms from. If no category is given, all terms from all categories will be listed.

> **Flags:**
> * -d, --done: (optional - bool) Use this flag to only list the terms that have been marked as 'done'.
> * -g, --categories: (optional - bool) Use this flag to list all categories, without any terms.

*Note: the -g flag must be used by itself, and cannot be used in conjunction with a category argument or any other flags.*


### Done
> *Mark a term as 'done'*

> **Arguments:**
> * category: (optional - string) The category that the term to mark as 'done' is within.

> **Flags:**
> * -t, --term: (optional - string) The term to mark as 'done.

### Remove
> *Remove a term from a category.*

> **Arguments:**
> * category: (optional - string) The category that the term to remove is within.

> **Flags:**
> * -t, --term: (optional - string) The term to remove from the category.

### Ask
> *Ask a question and have an answer generated by the gemini-1.5-pro model from Google.*

> **Arguments:**
> * question: (required - string) The question to ask the model.

> **Flags:**
> * -c, --category: (optional - string) The context/category of the question.
> * -f, --file: (optional - string) The file to pass the contents of to the prompt. You can pass only the name of the file without the entire path, and the tool will search for the file in all subdirectories of the current directory.
> * -b, --brief: (optional - bool) Use this flag to have the answer given in as few words as possible.
> * -v, --verbose: (optional - bool) Use this flag to get a lengthy, detailed response.

> *Note: The -b and -v flags are mutually exclusive, and cannot be used together.*

### Help
> *Display the help information for a command.*

> **Arguments:**
> * command: (optional - string) The command to display help information for.

> *Note: If no command is given, the help information for all commands will be displayed.*


Although the **add**, **list**, **done**, and **remove** commands don't require a category or term to be given, if the tool would need that information to complete the action requested, the user will be prompted to provide it through either a selection menu or manually typing it in.

---
## Example Usage
**Ask command with -b flag**
