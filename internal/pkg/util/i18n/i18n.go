// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package i18n

// https://github.com/kubernetes/kubectl/blob/release-1.24/pkg/util/i18n/i18n.go

import (
	"archive/zip"
	"bytes"
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/chai2010/gettext-go/gettext"
)

//go:embed translations
var translations embed.FS

var knownTranslations = map[string][]string{
	"q": {
		"default",
		"en_US",
		"zh_CN",
	},
}

func loadSystemLanguage() string {
	langStr := os.Getenv("LC_ALL")
	if langStr == "" {
		langStr = os.Getenv("LC_MESSAGES")
	}
	if langStr == "" {
		langStr = os.Getenv("LANG")
	}
	if langStr == "" {
		// log.Infof("Couldn't find the LC_ALL, LC_MESSAGES or LANG environment variables, defaulting to en_US")
		return "default"
	}
	pieces := strings.Split(langStr, ".")
	if len(pieces) != 2 {
		// log.Debugf("Unexpected system language (%s), defaulting to en_US", langStr)
		return "default"
	}
	return pieces[0]
}

func findLanguage(root string, getLanguageFn func() string) string {
	langStr := getLanguageFn()

	translations := knownTranslations[root]
	for ix := range translations {
		if translations[ix] == langStr {
			return langStr
		}
	}
	// log.Infof("Couldn't find translations for %s, using default", langStr)
	return "default"
}

func LoadTranslations(root string, getLanguageFn func() string) error {
	if getLanguageFn == nil {
		getLanguageFn = loadSystemLanguage
	}

	langStr := findLanguage(root, getLanguageFn)
	translationFiles := []string{
		fmt.Sprintf("%s/%s/LC_MESSAGES/q.po", root, langStr),
		fmt.Sprintf("%s/%s/LC_MESSAGES/q.mo", root, langStr),
	}

	// TODO: list the directory and load all files.
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// Make sure to check the error on Close.
	for _, file := range translationFiles {
		filename := "translations/" + file
		f, err := w.Create(file)
		if err != nil {
			return err
		}
		data, err := translations.ReadFile(filename)
		if err != nil {
			return err
		}
		if _, err := f.Write(data); err != nil {
			return nil
		}
	}
	if err := w.Close(); err != nil {
		return err
	}
	gettext.BindTextdomain("q", root+".zip", buf.Bytes())
	gettext.Textdomain("q")
	gettext.SetLocale(langStr)
	return nil
}

func T(defaultValue string, args ...int) string {
	if len(args) == 0 {
		return gettext.PGettext("", defaultValue)
	}
	return fmt.Sprintf(gettext.PNGettext("", defaultValue, defaultValue+".plural", args[0]),
		args[0])
}
