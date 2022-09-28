package appi18n

var multiLocales *MultiLocales

type MultiLocales struct {
	LocalesGroup *[]AppLocales
}

func GetMultiLocales() *MultiLocales {
	return multiLocales
}

func (m *MultiLocales) GetSpeLocales(lang string) AppLocales {
	for _, appLocales := range *m.LocalesGroup {
		if appLocales.Lang == lang {
			return appLocales
		}
	}
	return (*multiLocales.LocalesGroup)[0]
}
