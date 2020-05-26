package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// Image struct for images
type Image struct {
	ID    string `xml:"id,attr"`
	Width string `xml:"width,attr"`
	Height string `xml:"height,attr"`
	Boxes []Box  `xml:"box"`
}

// Box struct for each box
type Box struct {
	Label      string      `xml:"label,attr"`
	XTL        string      `xml:"xtl,attr"`
	YTL        string      `xml:"ytl,attr"`
	XBR        string      `xml:"xbr,attr"`
	YBR        string      `xml:"ybr,attr"`
	Attributes []Attribute `xml:"attribute"`
}

// Attribute struct for storing attribute data
type Attribute struct {
	Value string `xml:",chardata"`
	Name  string `xml:"name,attr"`
}

func classFinder(attributes *[]Attribute) (error, string) {
	var helmet string
	var mask string
	if length := len(*attributes); length == 2 {
		for _, item := range *attributes {
			if item.Name == "has_safety_helmet" {
				helmet = item.Value
			} else if item.Name == "mask" {
				mask = item.Value
			}
		}
		if helmet == "yes" {
			if mask == "yes" {
				return nil, "0"
			} else if mask == "no" {
				return nil, "1"
			} else if mask == "invisible" {
				return nil, "3"
			}
		} else if helmet == "no" {
			if mask == "yes" {
				return nil, "2"
			} else if mask == "no" {
				return nil, "5"
			} else if mask == "invisible" {
				return nil, "4"
			}
		}
	}
	return fmt.Errorf("doesn't fit requirement for classifying"), "-1"
}

func main() {
	input := "annotations.xml"
	output := "output/"

	// Opening our xml file

	xmlFile, err := os.Open(input)
	if err != nil {
		log.Fatalf("Unable to open the file: %v", err)
	}

	fmt.Println("Successfully Opened " + input)
	d := xml.NewDecoder(xmlFile)

	for {
		t, tokenErr := d.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				break
			} else {
				log.Printf("Unable read a token: %v", tokenErr)
			}
		}
		switch t := t.(type) {
		case xml.StartElement:
			if t.Name.Local == "image" {
				var image Image
				if err := d.DecodeElement(&image, &t); err != nil {
					log.Printf("Unable decode record: %v\n", err)
				}
				file, err := os.Create(output + image.ID + ".txt")
				if err != nil {
					log.Printf("Unable to create output file for %s \n Error: %v\n", image.ID, err)
				}
				var str []string
				for _, box := range image.Boxes {
					if box.Label == "head" {
						err, class := classFinder(&box.Attributes)
						if err != nil && class == "-1" {
							continue
						}
						h, _ := strconv.ParseFloat(image.Height, 64)
						w, _ := strconv.ParseFloat(image.Width, 64)
						height := 1/h
						width := 1/w
						XTL, _ := strconv.ParseFloat(box.XTL, 64)
						XBR, _ := strconv.ParseFloat(box.XBR, 64)
						YTL, _ := strconv.ParseFloat(box.YTL, 64)
						YBR, _ := strconv.ParseFloat(box.YBR, 64)
						x := fmt.Sprintf("%f", (XTL + XBR)/2 * width)
						y := fmt.Sprintf("%f" ,(YTL + YBR)/2 * height)
						N := fmt.Sprintf("%f", (XBR - XTL) * width)
						H := fmt.Sprintf("%f", (YBR - YTL) * height)
						metrics := class + " " + x + " " + y + " " + N + " " + H
						str = append(str, metrics)
					}
				}
				for _, v := range str {
					_, err := fmt.Fprintln(file, v)
					if err != nil {
						log.Printf("Unable to write into %s\n %v\n", image.ID+".txt", err)
					}
				}
				fileErr := file.Close()
				if fileErr != nil {
					log.Printf("Unable to close the image data file: %v", fileErr)
				} else {
					log.Printf("Image -%s- completed.", image.ID)
				}
			}
		}
	}

	err = xmlFile.Close()
	if err != nil {
		log.Fatalf("Unable to close the file: %v", err)
	}
}
