package main

import "strings"

func (e *Editor) checkContents() {
	// Check if the first line is special
	firstLine := e.Line(0)
	if strings.HasPrefix(firstLine, "#!") { // The line starts with a shebang
		words := strings.Split(firstLine, " ")
		lastWord := words[len(words)-1]
		if strings.Contains(lastWord, "/") {
			words = strings.Split(lastWord, "/")
			lastWord = words[len(words)-1]
		}
		switch lastWord {
		case "python":
			e.mode = modePython
		case "bash", "fish", "zsh", "tcsh", "ksh", "sh", "ash":
			e.mode = modeShell
		}
	} else if strings.HasPrefix(firstLine, "# $") {
		// Most likely a csh script on FreeBSD
		e.mode = modeShell
	} else if strings.HasPrefix(firstLine, "#") {
		e.firstLineHash = true
	}
	// If more lines start with "# " than "// " or "/* ", and mode is blank,
	// set the mode to modeConfig and enable syntax highlighting.
	if e.mode == modeBlank {
		hashComment := 0
		slashComment := 0
		for _, line := range strings.Split(e.String(), "\n") {
			if strings.HasPrefix(line, "# ") {
				hashComment++
			} else if strings.HasPrefix(line, "/") { // Count all lines starting with "/" as a comment, for this purpose
				slashComment++
			}
		}
		if hashComment > slashComment {
			e.mode = modeConfig
			e.syntaxHighlight = true
		}
	}
	// If the mode is modeOCaml and there are no ";;" strings, switch to Standard ML
	if e.mode == modeOCaml {
		if !strings.Contains(e.String(), ";;") {
			e.mode = modeStandardML
		}
	}
}

func (e *Editor) adjustTabsAndSpaces() {
	// Additional per-mode considerations, before launching the editor
	switch e.mode {
	case modeMakefile, modePython, modeCMake, modeJava, modeKotlin, modeZig, modeScala:
		e.tabs = TabsSpaces{4, false}
	case modeShell, modeConfig, modeHaskell, modeVim, modeLua, modeObjectPascal:
		e.tabs = TabsSpaces{2, false}
	case modeAda:
		e.tabs = TabsSpaces{3, false}
	case modeMarkdown, modeText, modeBlank:
		e.rainbowParenthesis = false
	}
}
