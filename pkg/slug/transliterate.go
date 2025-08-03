package slug

import (
	"strings"
	"unicode"
)

func Transliterate(text string) string {
	cyrillicToLatin := map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
		'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
		'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
		'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch", 'ъ': "",
		'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
		'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D", 'Е': "E", 'Ё': "Yo",
		'Ж': "Zh", 'З': "Z", 'И': "I", 'Й': "Y", 'К': "K", 'Л': "L", 'М': "M",
		'Н': "N", 'О': "O", 'П': "P", 'Р': "R", 'С': "S", 'Т': "T", 'У': "U",
		'Ф': "F", 'Х': "Kh", 'Ц': "Ts", 'Ч': "Ch", 'Ш': "Sh", 'Щ': "Sch", 'Ъ': "",
		'Ы': "Y", 'Ь': "", 'Э': "E", 'Ю': "Yu", 'Я': "Ya",

		'Ә': "A'", 'ә': "a'", 'Ғ': "Gh", 'ғ': "gh", 'Қ': "Q", 'қ': "q",
		'Ң': "Ng", 'ң': "ng", 'Ө': "O'", 'ө': "o'", 'Ұ': "U'", 'ұ': "u'",
		'Ү': "U'", 'ү': "u'", 'Һ': "H", 'һ': "h",
	}

	var result []rune
	for _, char := range text {
		if latin, found := cyrillicToLatin[char]; found {
			result = append(result, []rune(latin)...)
		} else if unicode.IsSpace(char) || char == '-' || char == '_' {
			result = append(result, '-')
		} else if unicode.IsLetter(char) || unicode.IsDigit(char) {
			result = append(result, char)
		}
	}
	resultStr := strings.Join(strings.Fields(string(result)), "-")

	resultStr = strings.ToLower(resultStr)

	return resultStr
}
