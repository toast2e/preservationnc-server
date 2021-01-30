package html

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/toast2e/preservationnc-server/reps"
	"golang.org/x/net/html"
)

// PropertyFinder represents an object that can find properties
type PropertyFinder interface {
	// FindProperties is used to find properties
	FindProperties() ([]reps.Property, error)
}

// Crawler is used to parse html
type Crawler struct {
	client http.Client
}

// NewCrawler returns a new crawler
func NewCrawler(c http.Client) Crawler {
	return Crawler{client: c}
}

// FindProperties parses properties from html
func (c *Crawler) FindProperties() ([]reps.Property, error) {
	resp, err := c.client.Get("https://www.presnc.org/property-listing/all-properties/")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response from server: %s", resp.Status)
	}
	tokenizer := html.NewTokenizer(resp.Body)
	properties := make([]reps.Property, 0)
	for {
		_, err := c.advanceNextToken(tokenizer)
		if err != nil {
			if err == io.EOF {
				//end of the file, break out of the loop
				break
			} else {
				return properties, err
			}
		}

		// find property tokens
		token := tokenizer.Token()
		if token.Data == "div" {
			isProp, id := c.containsProperty(token.Attr)
			if isProp {
				log.Printf("got property div token with id = %v %s %v %s", token.Type, token.Data, token.Attr, id)
				// found a property, next <a> tag should be a link to the details
				token, err := c.findStartTokenData("a", tokenizer)
				if err != nil {
					return properties, err
				}
				log.Printf("got anchor tag token for id = %v %s %v %s", token.Type, token.Data, token.Attr, id)
				// TODO one or more of the properties is formatted differently than the others because of course it is FIXIT
				prop, err := c.propertyFromLink(id, token.Attr[0].Val)
				if err != nil {
					return properties, err
				}
				properties = append(properties, prop)
			}
		}
	}
	return properties, nil
}

func (c *Crawler) containsProperty(attr []html.Attribute) (bool, string) {
	for _, a := range attr {
		if a.Key == "id" && strings.Contains(a.Val, "property-") {
			if a.Val != "property-info" {
				id := strings.Split(a.Val, "-")[1]
				return true, id
			}
			return false, ""
		}
	}
	return false, ""
}

func (c *Crawler) propertyFromLink(id string, url string) (reps.Property, error) {
	log.Printf("getting info for property %s from url %s", id, url)
	resp, err := c.client.Get(url)
	if err != nil {
		return reps.Property{}, err
	}
	if resp.StatusCode != 200 {
		return reps.Property{}, fmt.Errorf("unexpected response from server: %d", resp.StatusCode)
	}

	prop := reps.Property{ID: id}

	// find and set the name of the property which is in an <h1/> tag
	tokenizer := html.NewTokenizer(resp.Body)
	token, err := c.findTokenData("h1", tokenizer)
	if err != nil {
		return reps.Property{}, err
	}

	token, err = c.findTokenType(html.TextToken, tokenizer)
	if err != nil {
		return reps.Property{}, err
	}
	prop.Name = token.Data

	// find the rest of the property info which is in a <div/> tag with id=single-property-info
	token, err = c.findTokenWithAttributeValue("div", "id", "single-property-info", tokenizer)
	if err != nil {
		return reps.Property{}, err
	}

	// find the street address
	token, err = c.findTokenWithAttributeValue("span", "class", "street-address", tokenizer)
	if err != nil {
		return reps.Property{}, err
	}

	token, err = c.findTokenType(html.TextToken, tokenizer)
	if err != nil {
		return reps.Property{}, err
	}
	prop.Location.Address = strings.TrimSpace(token.Data)

	// find the city
	token, err = c.findTokenWithAttributeValue("span", "class", "locality", tokenizer)
	if err != nil {
		return reps.Property{}, err
	}

	token, err = c.findTokenType(html.TextToken, tokenizer)
	if err != nil {
		return reps.Property{}, err
	}
	prop.Location.City = strings.TrimSpace(token.Data)

	// find the state
	token, err = c.findTokenWithAttributeValue("span", "class", "region", tokenizer)
	if err != nil {
		return reps.Property{}, err
	}

	token, err = c.findTokenType(html.TextToken, tokenizer)
	if err != nil {
		return reps.Property{}, err
	}
	prop.Location.State = strings.TrimSpace(token.Data)

	// find the zip
	token, err = c.findTokenWithAttributeValue("span", "class", "postal-code", tokenizer)
	if err != nil {
		return reps.Property{}, err
	}

	token, err = c.findTokenType(html.TextToken, tokenizer)
	if err != nil {
		return reps.Property{}, err
	}
	prop.Location.Zip = strings.TrimSpace(token.Data)

	// find the county
	token, err = c.findTokenWithAttributeValue("span", "class", "county", tokenizer)
	if err != nil {
		return reps.Property{}, err
	}

	token, err = c.findTokenType(html.TextToken, tokenizer)
	if err != nil {
		return reps.Property{}, err
	}
	prop.Location.County = strings.TrimSpace(token.Data)

	// find the price which should be the next <li/> tag
	token, err = c.findStartTokenData("li", tokenizer)
	if err != nil {
		return reps.Property{}, err
	}

	token, err = c.findTokenType(html.TextToken, tokenizer)
	if err != nil {
		return reps.Property{}, err
	}
	priceString := strings.TrimSpace(token.Data)
	// only attempt to parse a price if the value starts with a '$'.
	// not all properties include a price
	if priceString[0] == '$' {
		// trim off the dollar sign
		priceString = priceString[1:]
		// remove any commas
		priceString = strings.ReplaceAll(priceString, ",", "")
		prop.Price, err = strconv.ParseFloat(priceString, 32)
		if err != nil {
			return reps.Property{}, err
		}
	}

	log.Printf("Parsed property from %s: %v", url, prop)
	return prop, nil

}

func (c *Crawler) containsAttributeWithValue(key string, value string, attr []html.Attribute) bool {
	for _, a := range attr {
		if a.Key == key && strings.Contains(a.Val, value) {
			return true
		}
	}
	return false
}

func (c *Crawler) findTokenType(tokenType html.TokenType, tokenizer *html.Tokenizer) (html.Token, error) {
	for {
		nextTokenType, err := c.advanceNextToken(tokenizer)
		if err != nil {
			return html.Token{}, err
		}

		if nextTokenType == tokenType {
			return tokenizer.Token(), nil
		}
	}
}

func (c *Crawler) findTokenData(data string, tokenizer *html.Tokenizer) (html.Token, error) {
	for {
		_, err := c.advanceNextToken(tokenizer)
		if err != nil {
			return html.Token{}, err
		}
		token := tokenizer.Token()

		if token.Data == data {
			return token, nil
		}
	}
}

func (c *Crawler) findStartTokenData(data string, tokenizer *html.Tokenizer) (html.Token, error) {
	for {
		_, err := c.advanceNextToken(tokenizer)
		if err != nil {
			return html.Token{}, err
		}
		token := tokenizer.Token()

		if token.Type == html.StartTagToken && token.Data == data {
			return token, nil
		}
	}
}

func (c *Crawler) findTokenWithAttributeValue(data string, attrKey string, attrValue string, tokenizer *html.Tokenizer) (html.Token, error) {
	for {
		token, err := c.findTokenData(data, tokenizer)
		if err != nil {
			return token, err
		}
		if c.containsAttributeWithValue(attrKey, attrValue, token.Attr) {
			return token, nil
		}
	}
}

func (c *Crawler) advanceNextToken(tokenizer *html.Tokenizer) (html.TokenType, error) {
	//get the next token type
	tokenType := tokenizer.Next()

	//if it's an error token, we either reached
	//the end of the file, or the HTML was malformed
	if tokenType == html.ErrorToken {
		err := tokenizer.Err()
		if err != nil {
			return tokenType, err
		}
	}

	return tokenType, nil
}
