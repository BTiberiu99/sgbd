package translate

import "sgbd4/go/legend"

var (
	message = map[string]string{
		legend.MessageConnectionSuccess:       "Conexiunea a fost createa cu succes",
		legend.MessageConnectionFail:          "Conexiunea a esuat cu eroarea :\n%s",
		legend.MessageConnectionSuccessRemove: "Conexiunea %s a fost stearsa cu succes",
		legend.MessageNoConnection:            "Trebuie sa stabilit intai o conexiune!",
		legend.MessageConnectionExist:         "Conexiunea %s exista deja !",
		legend.MessageConnectionNotExist:      "Conexiunea %s nu exista! Trebuie sa o creati intai!",
		legend.MessageSuccessAddNotNULL:       "S-a adaugat cu succes constrangerea not null pentru coloana %s",
		legend.MessageFailAddNotNULL:          "Nu s-a putu adauga constrangearae not null pentru coloana %s",
		legend.MessagePrimaryKeySuccess:       "Cheia primara %s a fost adauga cu succes",
		legend.MessagePrimaryKeyFail:          "Cheia primara %s nu a putut fii adaugata",
		legend.MessagePrimaryKeySuccessFix:    "Cheia primara a fost reparata cu succes",
		legend.MessagePrimaryKeyFailFix:       "Cheia primara nu a putut fii reparata, error:%s ",
	}
)
