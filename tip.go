package main

import (
	"fmt"
	"log"
	"os"
)

const MaxEntities = 5000

//go:generate stringer -type=ListKind
type ListKind int

const (
	Files ListKind = iota
	Mails
	URLs
	Detections
)

type Tip struct {
	lists [Detections + 1][]string
}

func NewTip() *Tip {
	return &Tip{}
}

func (t *Tip) Add(list ListKind, file string) {
	t.lists[list] = append(t.lists[list], file)
}

func (t *Tip) Count() (result int) {
	for i := Files; i <= Detections; i++ {
		result += len(t.lists[i])
	}
	return
}

type TipList struct {
	tips []Tip
}

func NewTipList() *TipList {
	return &TipList{
		tips: []Tip{},
	}
}

func (t *TipList) Add(list ListKind, file string) {
	if len(t.tips) == 0 {
		t.tips = append(t.tips, Tip{})
	}
	if t.tips[len(t.tips)-1].Count() == MaxEntities {
		t.tips = append(t.tips, Tip{})
	}
	t.tips[len(t.tips)-1].Add(list, file)
}

func (t *TipList) GenerateResult(baseFileName string) error {
	for i, tip := range t.tips {
		fileName := fmt.Sprintf("%s%04d.txt", baseFileName, i)
		log.Printf("Output file: %s", fileName)
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		for kind, list := range tip.lists {
			if len(list) == 0 {
				continue
			}
			fmt.Fprintf(file, "[%v]\n", ListKind(kind))
			for _, each := range list {
				if _, err = fmt.Fprintf(file, "%s\n", each); err != nil {
					return err
				}
			}
		}
		if err := file.Close(); err != nil {
			return err
		}
	}
	return nil
}
