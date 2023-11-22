package dao

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
	"uttc-hackathon/database"
	"uttc-hackathon/model"

	"github.com/go-sql-driver/mysql"
	"github.com/ikawaha/kagome/tokenizer"
)

type RelatedItemDAO interface {
	GetItemsSimilarity(itemID int, itemCategoriesName string) ([]model.Item, error)
}

type relatedItemDAO struct {
	db *sql.DB
}

func NewRelatedItemDAO(db *sql.DB) RelatedItemDAO {
	return &relatedItemDAO{
		db: db,
	}
}

func (rid *relatedItemDAO) GetItemsSimilarity(itemID int, itemCategoriesName string) ([]model.Item, error) {
	var main_item model.Item
	var CreatedAt, UpdatedAt mysql.NullTime

	query := fmt.Sprintf(
		"SELECT i.id, i.user_firebase_uid, i.title, i.author, i.link, i.explanation, i.likes, i.item_categories_id, icat.name AS item_category_name, i.created_at, i.updated_at, CONCAT(i.title, ' ', i.author, ' ', i.link, ' ', i.explanation) AS concatenated_text FROM %s AS i "+
			"INNER JOIN item_curriculums AS ic ON i.id = ic.item_id AND i.item_categories_id = ic.item_categories_id "+
			"INNER JOIN item_categories AS icat ON ic.item_categories_id = icat.id "+
			"WHERE i.id = ? ", itemCategoriesName)
	err := database.DB.QueryRow(query, itemID).Scan(
		&main_item.ID, &main_item.UserFirebaseUID, &main_item.Title, &main_item.Author, &main_item.Link, &main_item.Explanation,
		&main_item.Likes, &main_item.ItemCategoriesID, &main_item.ItemCategoriesName, &CreatedAt, &UpdatedAt, &main_item.ConcatenatedText)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, err
	}

	main_item.CreatedAt = CreatedAt.Time
	main_item.UpdatedAt = UpdatedAt.Time

	var items []model.Item
	items, err = getALLItems()
	fmt.Println("items: ", items)
	if err != nil {
		return nil, err
	}

	for index, item := range items {
		similarity := CosineSimilarity(main_item.ConcatenatedText, item.ConcatenatedText, "ja")
		items[index].Similarity = similarity
		fmt.Println("similarity: ", item.Similarity)
	}

	return items, nil
}

// CosineSimilarity calculates the cosine similarity between two text strings.
func CosineSimilarity(text1, text2 string, language string) float64 {
	// Tokenize and remove stopwords from the input text
	var words1, words2 []string
	fmt.Println("text1: ", text1)
	fmt.Println("text2: ", text2)

	if language == "ja" {
		// For Japanese text, use custom tokenization logic
		words1 = tokenizeJapanese(text1)
		words2 = tokenizeJapanese(text2)
	} else {
		// For other languages (including English), use basic tokenization and stopword removal
		words1 = removeStopwords(tokenize(text1))
		words2 = removeStopwords(tokenize(text2))
	}
	fmt.Println(words1)
	fmt.Println(words2)

	// Create a map to store word frequencies for both texts
	freqMap1 := make(map[string]int)
	freqMap2 := make(map[string]int)

	// Calculate word frequencies for text1
	for _, word := range words1 {
		freqMap1[word]++
	}

	// Calculate word frequencies for text2
	for _, word := range words2 {
		freqMap2[word]++
	}

	fmt.Println("freaMap1: ", freqMap1)
	fmt.Println("freqMap2: ", freqMap2)

	// Calculate the dot product of word vectors
	dotProduct := 0.0
	for word, freq1 := range freqMap1 {
		if freq2, found := freqMap2[word]; found {
			dotProduct += float64(freq1 * freq2)
		}
	}
	fmt.Println("dotProduct", dotProduct)

	// Calculate the magnitude of each word vector
	magnitude1 := 0.0
	for _, freq1 := range freqMap1 {
		magnitude1 += float64(freq1 * freq1)
	}
	magnitude1 = math.Sqrt(magnitude1)

	magnitude2 := 0.0
	for _, freq2 := range freqMap2 {
		magnitude2 += float64(freq2 * freq2)
	}
	magnitude2 = math.Sqrt(magnitude2)
	fmt.Println(magnitude1)
	fmt.Println(magnitude2)

	// Calculate the cosine similarity
	if magnitude1 != 0 && magnitude2 != 0 {
		similarity := dotProduct / (magnitude1 * magnitude2)
		fmt.Println("similarity ", similarity)
		return similarity
	}

	return 0.0 // Default similarity when one or both vectors are zero
}

// Tokenize text into words
func tokenize(text string) []string {
	return strings.Fields(text)
}

// Tokenize Japanese text using custom logic
func tokenizeJapanese(text string) []string {
	t := tokenizer.New()
	tokens := t.Tokenize(text)

	var result []string
	for _, token := range tokens {
		// Ignore punctuation and symbols
		if token.Class == tokenizer.DUMMY {
			continue
		}
		result = append(result, token.Surface)
	}

	return result
}

// Remove stopwords from a list of words (Japanese specific)
func removeStopwords(words []string) []string {
	// Use a list of Japanese stopwords
	stopWords := []string{"の", "に", "は", "を", "た", "が", "で", "て", "と", "し", "れ", "さ", "ある", "いる", "も", "る", "な", "ない", "など", "なら", "か", "から", "こと", "この", "これ", "それ", "だ", "だけ", "でも", "って", "という", "まで", "もの", "や", "よ", "ら", "られ", "れる", "わ", "を通じ", "とともに", "ただし", "について", "のみ", "それぞれ", "または", "いう", "ば", "ながら", "へ", "より"} // This is an example list. Replace it with your own list

	// Create a map for efficient lookup
	stopWordsMap := make(map[string]bool)
	for _, word := range stopWords {
		stopWordsMap[word] = true
	}

	// Remove stopwords from the list of words
	filteredWords := []string{}
	for _, word := range words {
		if !stopWordsMap[word] {
			filteredWords = append(filteredWords, word)
		}
	}

	return filteredWords
}

func getALLItems() ([]model.Item, error) {
	var items []model.Item

	// Get all curriculum IDs from the categories table
	curriculumIDsAll, err := getAllCurriculumIDs()
	if err != nil {
		return nil, err
	}

	// Create a map to track unique item categories for each item
	uniqueItemCategories := make(map[string]bool)

	// Iterate through curriculum IDs and table names to retrieve filtered items
	for _, curriculumID := range curriculumIDsAll {
		// Iterate through the tableNames (or categories) and retrieve items for each curriculum
		tableNames, err := getTableNames()
		if err != nil {
			return nil, err
		}

		// Iterate through table names to retrieve items for each category
		for _, tableName := range tableNames {
			// Build the query based on curriculum and category
			query := fmt.Sprintf("SELECT i.id, i.user_firebase_uid, i.title, i.author, i.link, i.explanation, i.likes, i.item_categories_id, icat.name AS item_category_name, i.created_at, i.updated_at, CONCAT(i.title, ' ', i.author, ' ', i.link, ' ', i.explanation) AS concatenated_info FROM %s AS i "+
				"INNER JOIN item_curriculums AS ic ON i.id = ic.item_id AND i.item_categories_id = ic.item_categories_id "+
				"INNER JOIN item_categories AS icat ON ic.item_categories_id = icat.id "+
				"WHERE ic.curriculum_id = ? ", tableName)

			rows, err := database.DB.Query(query, curriculumID) // Provide both curriculumID and tableName as arguments

			if err != nil {
				// Log the error for debugging
				fmt.Printf("Error executing query: %v", err)
				return nil, err
			}
			defer rows.Close()

			// Retrieve and append data from the table
			for rows.Next() {
				var item model.Item

				var CreatedAt, UpdatedAt mysql.NullTime

				// Scan the data into the item struct
				if err := rows.Scan(
					&item.ID, &item.UserFirebaseUID, &item.Title, &item.Author, &item.Link, &item.Explanation,
					&item.Likes, &item.ItemCategoriesID, &item.ItemCategoriesName, &CreatedAt, &UpdatedAt, &item.ConcatenatedText); err != nil {
					fmt.Println("error in scan")
					fmt.Println(err)
					return nil, err
				}

				// Convert NullTime to time.Time if not null
				if CreatedAt.Valid {
					item.CreatedAt = CreatedAt.Time
				}
				if UpdatedAt.Valid {
					item.UpdatedAt = UpdatedAt.Time
				}

				// Retrieve and add the curriculum IDs for this item
				curriculumIDs, err := getCurriculumIDsForItem(item.ID, item.ItemCategoriesID)
				if err != nil {
					fmt.Println("error: Failed to retrieve curriculum IDs for item")
					return nil, err
				}
				item.CurriculumIDs = curriculumIDs // Assign the curriculum IDs to the item

				// Create a unique identifier for the item based on item_id and item_categories_id
				itemIdentifier := fmt.Sprintf("%d-%d", item.ID, item.ItemCategoriesID)

				// Check if this item is unique based on the item identifier
				if !uniqueItemCategories[itemIdentifier] {
					// If it's unique, append it to the items slice and mark it as displayed
					items = append(items, item)
					uniqueItemCategories[itemIdentifier] = true
				}

			}
		}
	}
	return items, nil

}
