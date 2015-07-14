package file

import (
	"errors"
	"os"
	"strconv"
	"unicode"
)

// ParseMode parses a string s as the parameter of chmod command. It returns the
// FileMode which is applied the parsed command into the base FileMode.
//
// Chmod has following 2 types of parameter, absolute and symbolic mode.
//
// Absolute mode:
//
//   Absolute mode assumes the string s as an octal number of the file mode.
//
// Symbolic mode:
//
//   The symbolic mode is described by the following grammar.
//
//   mode   ::= clause [, clause ...]
//   clause ::= [who ...] [action ...] action
//   action ::= op [perm ...]
//   who    ::= a | u | g | o
//   op     ::= + | - | =
//   perm   ::= r | s | t | w | x | X | u | g | o
//
func ParseMode(s string, base os.FileMode) (os.FileMode, error) {
	if unicode.IsDigit(rune(s[0])) {
		return parseModeDigits(s)
	}
	return parseModeSymbolic(s, base)
}

func parseModeDigits(s string) (os.FileMode, error) {
	i, err := strconv.ParseInt(s, 8, 32)
	if err != nil {
		return 0, err
	}
	return os.FileMode(i), nil
}

func parseModeSymbolic(s string, base os.FileMode) (os.FileMode, error) {
	// log.Println("-- " + s)
	perm := uint32(base.Perm())
	sbitsMask := uint32(os.ModeSetuid | os.ModeSetgid | os.ModeSticky)
	sbits := uint32(base) & sbitsMask

	for i := 0; i < len(s); {
		// user, group, other, or all
		var whom uint32

		for j, c := range s[i:] {
			// log.Printf("whom: %d %s\n", i, string(s[i]))
			if c == 'u' {
				whom = whom | 0700
			} else if c == 'g' {
				whom = whom | 0070
			} else if c == 'o' {
				whom = whom | 0007
			} else if c == 'a' {
				whom = whom | 0777
			} else {
				if j == 0 {
					whom = 0777
				}
				break
			}
			i++
		}
		if i >= len(s) {
			return 0, errors.New("failed to parse the mode, operator is not found")
		}

		for i < len(s) {
			// operator
			// log.Printf("op: %d %s\n", i, string(s[i]))
			op := 0
			switch s[i] {
			case '-':
				op = '-'
			case '+':
				op = '+'
			case '=':
				op = '='
			default:
				return 0, errors.New("failed to parse the mode, invalid operator")
			}
			i++

			// permission bits to modify
			var smod uint32
			var mod uint32
			if i < len(s) {
				// log.Printf("mod: %d %s\n", i, string(s[i]))
				switch s[i] {
				case 'u':
					mod = (perm & 0700) >> 6
					i++
				case 'g':
					mod = (perm & 0070) >> 3
					i++
				case 'o':
					mod = (perm & 0007) >> 0
					i++
				default:
					for i < len(s) {
						// log.Printf("mod-f: %d %s\n", i, string(s[i]))
						if s[i] == 'r' {
							mod = mod | 4
						} else if s[i] == 'w' {
							mod = mod | 2
						} else if s[i] == 'x' {
							mod = mod | 1
						} else if s[i] == 'X' {
							if base.IsDir() {
								mod = mod | 1
							}
						} else if s[i] == 's' {
							if (whom & 0700) > 0 {
								smod = smod | uint32(os.ModeSetuid)
							}
							if (whom & 0070) > 0 {
								smod = smod | uint32(os.ModeSetgid)
							}
						} else if s[i] == 't' {
							smod = smod | uint32(os.ModeSticky)
						} else {
							break
						}
						i++
					}
				}
				mod = mod<<6 | mod<<3 | mod<<0
			}
			// if i < len(s) {
			//	log.Printf("mod-e: %d %s\n", i, string(s[i]))
			//}

			switch op {
			case '-':
				perm = perm - (whom & mod)
				sbits = sbits - smod
			case '+':
				perm = perm | (whom & mod)
				sbits = sbits | smod
			case '=':
				perm = (perm & (0777 ^ whom)) | (whom & mod)
				sbits = smod
			}

			if i < len(s) {
				// log.Printf(",: %d %s\n", i, string(s[i]))
				if s[i] == ',' {
					i++
					break
				}
			}
		}
	}

	mode := uint32(base)&((1<<32-1)-0777-sbitsMask) | sbits&sbitsMask | perm&0777
	return os.FileMode(mode), nil
}
