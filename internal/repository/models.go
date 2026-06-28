package repository

import "encoding/json"

type Item struct {
	ID       string     `json:"id"`
	Category string     `json:"category"`
	Name     TextObject `json:"name"`
	Color    string     `json:"color"`
	Status   struct {
		State string `json:"state"`
	} `json:"status"`
	InfoBlocks []InfoBlock `json:"infoBlocks"`
}

type TextObject struct {
	Type string
	Data any
}

type TranslationObject struct {
	Type  string            `json:"type"`
	Key   string            `json:"key"`
	Args  map[string]any    `json:"args"`
	Lines map[string]string `json:"lines"`
}

type SimpleTextObject struct {
	Text string `json:"text"`
}

func (to *TextObject) UnmarshalJSON(data []byte) error {
	textType, err := findType(data)
	if err != nil {
		return err
	}
	to.Type = textType

	switch to.Type {
	case "translation":
		var translation TranslationObject
		if err := json.Unmarshal(data, &translation); err != nil {
			return err
		}
		to.Data = translation

	case "text":
		var text SimpleTextObject
		if err := json.Unmarshal(data, &text); err != nil {
			return err
		}
		to.Data = text
	}

	return nil
}

type InfoBlock struct {
	Type string
	Data any
}

type TextBlock struct {
	Title TextObject `json:"title"`
	Text  TextObject `json:"text"`
}

type DamageBlock struct {
	StartDamage         float32 `json:"startDamage"`
	DamageDecreaseStart float32 `json:"damageDecreaseStart"`
	DamageDecreaseEnd   float32 `json:"damageDecreaseEnd"`
	EndDamage           float32 `json:"endDamage"`
	MaxDistance         float32 `json:"maxDistance"`
}

type ListBlock struct {
	Title    TextObject `json:"title"`
	Elements []Element  `json:"elements"`
}

type Element struct {
	Type string
	Data any
}

type NumericElement struct {
	Name  TextObject `json:"name"`
	Value float32    `json:"value"`
}

type KeyValueElement struct {
	Key   TextObject `json:"key"`
	Value TextObject `json:"value"`
}

type RangeElement struct {
	Name TextObject `json:"name"`
	Min  float32    `json:"min"`
	Max  float32    `json:"max"`
}

type TextElement struct {
	Text TextObject `json:"text"`
}

type ItemElement struct {
	Name TextObject `json:"name"`
}

func (e *Element) UnmarshalJSON(data []byte) error {
	elementType, err := findType(data)
	if err != nil {
		return err
	}
	e.Type = elementType

	switch e.Type {
	case "numeric":
		var numericElement NumericElement
		if err := json.Unmarshal(data, &numericElement); err != nil {
			return err
		}
		e.Data = numericElement
	case "key-value":
		var keyValueElement KeyValueElement
		if err := json.Unmarshal(data, &keyValueElement); err != nil {
			return err
		}
		e.Data = keyValueElement
	case "range":
		var rangeElement RangeElement
		if err := json.Unmarshal(data, &rangeElement); err != nil {
			return err
		}
		e.Data = rangeElement
	case "text":
		var textElement TextElement
		if err := json.Unmarshal(data, &textElement); err != nil {
			return err
		}
		e.Data = textElement
	case "item":
		var itemElement ItemElement
		if err := json.Unmarshal(data, &itemElement); err != nil {
			return err
		}
		e.Data = itemElement
	}

	return nil
}

func (ib *InfoBlock) UnmarshalJSON(data []byte) error {
	blockType, err := findType(data)
	if err != nil {
		return err
	}
	ib.Type = blockType

	switch ib.Type {
	case "text":
		var textBlock TextBlock
		if err := json.Unmarshal(data, &textBlock); err != nil {
			return err
		}
		ib.Data = textBlock

	case "damage":
		var damageBlock DamageBlock
		if err := json.Unmarshal(data, &damageBlock); err != nil {
			return err
		}
		ib.Data = damageBlock

	case "list":
		var listBlock ListBlock
		if err := json.Unmarshal(data, &listBlock); err != nil {
			return err
		}
		ib.Data = listBlock
	}

	return nil
}

func findType(data []byte) (string, error) {
	var typeFinder struct {
		Type string
	}
	if err := json.Unmarshal(data, &typeFinder); err != nil {
		return "", err
	}
	return typeFinder.Type, nil
}
