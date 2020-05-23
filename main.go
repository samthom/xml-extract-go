package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

// Annotations struct for entire file
type Annotations struct {
	XMLName xml.Name `xml:"annotations"`
	Images []Image `xml:"image"`
}

// Image struct for images
type Image struct {
	ID string `xml:"id,attr"`
	Boxes []Box `xml:"box"`
}

// Box struct for each box
type Box struct {
	Label string `xml:"label,attr"`
	XTL string `xml:"xtl,attr"`
	YTL string `xml:"ytl,attr"`
	XBR string `xml:"xbr,attr"`
	YBR string `xml:"ybr,attr"`
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
				//go func() {
				//
				//}()
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
						metrics := box.XTL +" "+ box.YTL +" "+ box.XBR +" "+ box.YBR
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