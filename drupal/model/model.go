// Provides a rudimentary model for JSON API resources, and a second model used for test expectations.
//
// Clients of this package are expected to formulate queries to Drupal's JSON API using JsonApiUrl, and unmarshalling
// the response to the JSON API specific structs, e.g. JsonApiIslandoraObj. Assertions may be made of the JSON API
// structs using the so-called 'Expected' structs, e.g. the ExpectedRepoObj.
//
// The 'Expected' structs are simple JSON objects whose values (e.g. title, name, or description) are predetermined
// by known data that were migrated into Drupal (e.g. via a migration) for the purpose of testing.  To assert whether
// an expected repository object was actually ingested, the JsonApiIslandoraObj would be compared to the ExpectedRepoObj.
package model

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jhu-idc/idc-golang/drupal/env"
	"github.com/jhu-idc/idc-golang/drupal/jsonapi"
)

const (
	// Constant for the Drupal node entity type
	Node = "node"
	// Constant for the IDC-specific collection entity type
	Collection = "collection_object"
	// Constant for the Islandora-specific repository object entity type
	RepositoryObject = "islandora_object"
	// Constant for the Image media bundle
	Image = "image"
	// Constant for the Document media bundle
	Document = "document"
	// Constant for the Video media bundle
	Video = "video"
	// Constant for the Audio media bundle
	Audio = "audio"
	// Constant for the Extracted Text media bundle
	ExtractedText = "extracted_text"
	// Constant for the File media bundle
	File = "file"
	// Constant for the FITS technical metadata media bundle
	Fits = "fits_technical_metadata"
	// Constant for the Remote Video media bundle
	RemoteVideo = "remote_video"
)

// Minimally models the elements present in a JSON API data element
type JsonApiData struct {
	// The Drupal type of the resource: a tuple of entity type and bundle
	Type jsonapi.DrupalType
	// The identifier of the resource contained in the data element, typically a UUID provided by Drupal
	Id string
}

// Resolve the reference of the data object, useful for references appearing within JSON API `relationships`.  This
// function formulates a JSON API query based on the type, bundle, and unique identifier of the object, and returns
// exactly one resource.
func (jad *JsonApiData) Resolve(t *testing.T, v interface{}) {
	u := jsonapi.JsonApiUrl{
		T:            t,
		// TODO FIXME the BaseUrl won't work as expected. Really the caller wants the BaseUrl that was used to retrieve
		//   the JsonApiData, which means we really need access to the JSON API 'links' object and use the 'self' href.
		//   But we can't do that easily right now.
		BaseUrl:      env.BaseUrlOr("https://islandora-idc.traefik.me"),
		DrupalEntity: jad.Type.Entity(),
		DrupalBundle: jad.Type.Bundle(),
		Filter:       "id",
		Value:        jad.Id,
	}

	u.GetSingle(v)
}

// ResolveWithBasicAuth behaves as Resolve, but issues the request with HTTP Basic Auth, using the supplied username and
// password
func (jad *JsonApiData) ResolveWithBasicAuth(t *testing.T, v interface{}, username string, password string) {
	u := jsonapi.JsonApiUrl{
		T:            t,
		// TODO FIXME the BaseUrl won't work as expected. Really the caller wants the BaseUrl that was used to retrieve
		//   the JsonApiData, which means we really need access to the JSON API 'links' object and use the 'self' href.
		//   But we can't do that easily right now.
		BaseUrl:      env.BaseUrlOr("https://islandora-idc.traefik.me"),
		DrupalEntity: jad.Type.Entity(),
		DrupalBundle: jad.Type.Bundle(),
		Filter:       "id",
		Value:        jad.Id,
		Username:     username,
		Password:     password,
	}

	u.GetSingle(v)
}

// Represents the results of a JSONAPI query for a single Person from the Person Taxonomy
type JsonApiPerson struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Name        string   `json:"name"`
			Dates       []string `json:"field_date"`
			Description struct {
				Value     string
				Format    string
				Processed string
			}
			PrimaryPartOfName       string   `json:"field_primary_part_of_name"`
			PreferredNamePrefix     []string `json:"field_preferred_name_prefix"`
			PreferredNameRest       []string `json:"field_preferred_name_rest"`
			PreferredNameSuffix     []string `json:"field_preferred_name_suffix"`
			PreferredNameFullerForm []string `json:"field_preferred_name_fuller_form"`
			PreferredNameNumber     []string `json:"field_preferred_name_number"`
			PersonAlternateName     []string `json:"field_person_alternate_name"`
			Authority               []struct {
				Uri    string
				Title  string
				Source string
			} `json:"field_authority_link"`
		} `json:"attributes"`
		JsonApiRelationships struct {
			Relationships struct {
				Data []struct {
					JsonApiData
					Meta map[string]string
				}
			} `json:"field_relationships"`
		} `json:"relationships"`
	} `json:"data"`
}

// Represents the results of a JSONAPI query for a single Access Rights Taxonomy Term
type JsonApiAccessRights struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Name        string
			Description struct {
				Value     string
				Format    string
				Processed string
			}
			Authority []struct {
				Uri    string
				Title  string
				Source string
			} `json:"field_authority_link"`
		} `json:"attributes"`
	} `json:"data"`
}

// Represents the results of a JSONAPI query for a single Islandora Access Taxonomy Term
type JsonApiIslandoraAccessTerms struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Name        string
			Description struct {
				Value     string
				Format    string
				Processed string
			}
		} `json:"attributes"`
		JsonApiRelationships struct {
			AccessTerms struct {
				Data []JsonApiData
			} `json:"parent"`
		} `json:"relationships"`
	} `json:"data"`
}

// Represents the results of a JSONAPI query for a single Copyright and Use Taxonomy Term
type JsonApiCopyrightAndUse struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Name        string
			Description struct {
				Value     string
				Format    string
				Processed string
			}
			Authority []struct {
				Uri    string
				Title  string
				Source string
			} `json:"field_authority_link"`
		} `json:"attributes"`
	} `json:"data"`
}

// Represents the results of a JSONAPI query for a single Family Taxonomy Term
type JsonApiFamily struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Name        string
			Date        []string `json:"field_date"`
			FamilyName  string   `json:"field_family_name"`
			Title       string   `json:"field_title_and_other_words"`
			Description struct {
				Value     string
				Format    string
				Processed string
			}
			Authority []struct {
				Uri    string
				Title  string
				Source string
			} `json:"field_authority_link"`
		} `json:"attributes"`
		JsonApiRelationships struct {
			Relationships struct {
				Data []struct {
					JsonApiData
					Meta map[string]string
				}
			} `json:"field_relationships"`
		} `json:"relationships"`
	} `json:"data"`
}

// Represents the results of a JSONAPI query for a single collection entity
type JsonApiCollection struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Title       string
			Description struct {
				Value    string
				LangCode string
			}
			ContactEmail     string   `json:"field_collection_contact_email"`
			ContactName      string   `json:"field_collection_contact_name"`
			CollectionNumber []string `json:"field_collection_number"`
			FindingAid       []struct {
				Uri   string
				Title string
			} `json:"field_finding_aid"`
		} `json:"attributes"`
		JsonApiRelationships struct {
			AltTitle struct {
				Data  []JsonApiLanguageValue
				Links struct {
					Related struct {
						Href string
					}
				}
			} `json:"field_alternative_title"`
			TitleLanguage struct {
				Data  JsonApiLanguageValue
				Links struct {
					Related struct {
						Href string
					}
				}
			} `json:"field_title_language"`
			Description struct {
				Data []JsonApiLanguageValue
			} `json:"field_description"`
			AccessTerms struct {
				Data []JsonApiData
			} `json:"field_access_terms"`
			MemberOf struct {
				Data JsonApiData
			} `json:"field_member_of"`
		} `json:"relationships"`
	} `json:"data"`
}

// Represents the results of a JSONAPI query for a single islandora object
type JsonApiIslandoraObj struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Title             string
			CollectionNumber  []string `json:"field_collection_number"`
			DateAvailable     string   `json:"field_date_available"`
			DateCopyrighted   []string `json:"field_date_copyrighted"`
			DateCreated       []string `json:"field_date_created"`
			DatePublished     []string `json:"field_date_published"`
			DigitalIdentifier []string `json:"field_digital_identifier"`
			DspaceIdentifier  struct {
				Uri   string
				Title string
			} `json:"field_dspace_identifier"`
			DspaceItemid string `json:"field_dspace_item_id"`
			Description  string
			Extent       []string `json:"field_extent"`
			FeaturedItem bool     `json:"field_featured_item"`
			FindingAid   []struct {
				Uri   string
				Title string
			} `json:"field_finding_aid"`
			GeoportalLink struct {
				Uri   string
				Title string
			} `json:"field_geoportal_link"`
			// TODO
			IsPartOf struct {
				Uri string
			} `json:"field_is_part_of"`
			Issn        string   `json:"field_issn"`
			ItemBarcode []string `json:"field_item_barcode"`
			JhirUri     struct {
				Uri   string
				Title string
			} `json:"field_jhir"`
			LibraryCatalogLink []struct {
				Uri   string
				Title string
			} `json:"field_library_catalog_link"`
			OclcNumber []string `json:"field_oclc_number"`
		} `json:"attributes"`
		JsonApiRelationships struct {
			Abstract struct {
				Data []JsonApiLanguageValue
			} `json:"field_abstract"`
			AccessRights struct {
				Data []JsonApiData
			} `json:"field_access_rights"`
			AccessTerms struct {
				Data []JsonApiData
			} `json:"field_access_terms"`
			AltTitle struct {
				Data []JsonApiLanguageValue
			} `json:"field_alternative_title"`
			Contributor struct {
				Data []RelData
			} `json:"field_contributor"`
			CopyrightAndUse struct {
				Data JsonApiData
			} `json:"field_copyright_and_use"`
			CopyrightHolder struct {
				Data []JsonApiData
			} `json:"field_copyright_holder"`
			Creator struct {
				Data []RelData
			} `json:"field_creator"`
			CustodialHistory struct {
				Data []JsonApiLanguageValue
			} `json:"field_custodial_history"`
			Description struct {
				Data []JsonApiLanguageValue
			} `json:"field_description"`
			DigitalPublisher struct {
				Data []JsonApiData
			} `json:"field_digital_publisher"`
			Genre struct {
				Data []JsonApiData
			} `json:"field_genre"`
			Language struct {
				Data []JsonApiData
			}
			Model struct {
				Data JsonApiData
			} `json:"field_model"`
			MemberOf struct {
				Data JsonApiData
			} `json:"field_member_of"`
			Publisher struct {
				Data []JsonApiData
			} `json:"field_publisher"`
			PublisherCountry struct {
				Data []JsonApiData
			} `json:"field_publisher_country"`
			ResourceType struct {
				Data []JsonApiData
			} `json:"field_resource_type"`
			SpatialCoverage struct {
				Data []JsonApiData
			} `json:"field_spatial_coverage"`
			Subject struct {
				Data []JsonApiData
			} `json:"field_subject"`
			TableOfContents struct {
				Data []JsonApiLanguageValue
			} `json:"field_table_of_contents"`
			TitleLanguage struct {
				Data JsonApiData
			} `json:"field_title_language"`
			DisplayHint struct {
				Data JsonApiData
			} `json:"field_display_hints"`
		} `json:"relationships"`
	} `json:"data"`
}

// Represents the results of a JSONAPI query for a single Genre Term
type JsonApiGenre struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Name        string
			Description struct {
				Value     string
				Format    string
				Processed string
			}
			Authority []struct {
				Uri    string
				Title  string
				Source string
			} `json:"field_authority_link"`
		} `json:"attributes"`
	} `json:"data"`
}

// Represents the results of a JSONAPI query for a single Geolocation Term
type JsonApiGeolocation struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Name    string
			Broader []struct {
				Uri   string
				Title string
			} `json:"field_broader"`
			GeoAltName  []string `json:"field_geo_alt_name"`
			Description struct {
				Value     string
				Format    string
				Processed string
			}
			Authority []struct {
				Uri    string
				Title  string
				Source string
			} `json:"field_authority_link"`
		} `json:"attributes"`
	} `json:"data"`
}

// Represents the results of a JSONAPI query for a single Resource Types Taxonomy Term
type JsonApiResourceType struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Name        string
			Description struct {
				Value     string
				Format    string
				Processed string
			}
			Authority []struct {
				Uri    string
				Title  string
				Source string
			} `json:"field_authority_link"`
		} `json:"attributes"`
	} `json:"data"`
}

// Represents the results of a JSONAPI query for a single Subject Term
type JsonApiSubject struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Name        string
			Description struct {
				Value     string
				Format    string
				Processed string
			}
			Authority []struct {
				Uri    string
				Title  string
				Source string
			} `json:"field_authority_link"`
		} `json:"attributes"`
	} `json:"data"`
}

// Represents the results of a JSONAPI query for a single Language Taxonomy Term
type JsonApiLanguage struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Name         string
			LanguageCode string `json:"field_language_code"`
			Description  struct {
				Value     string
				Format    string
				Processed string
			}
			Authority []struct {
				Uri    string
				Title  string
				Source string
			} `json:"field_authority_link"`
		} `json:"attributes"`
	} `json:"data"`
}

// Represents an element of a JSONAPI response that encapsulates a string value and a language taxonomy entity
//
// In the following example, the objects with a type `taxonomy_term--language` are represented by this struct.
//   "field_alternative_title": {
//    "data": [
//      {
//        "type": "taxonomy_term--language",
//        "id": "7397e0c4-df0a-4800-95af-afccc6ff64a5",
//        "meta": {
//          "value": "Moonrise Over Hernandez"
//        }
//      },
//      {
//        "type": "taxonomy_term--language",
//        "id": "bacfc5b6-b4b9-4239-8744-46dca6a91f0e",
//        "meta": {
//          "value": "Salida de la luna sobre Hernández"
//        }
//      }
//    ],
//    "links": {
//      "related": {
//        "href": "http://islandora-idc.traefik.me/jsonapi/node/islandora_object/815a4c04-0be5-44f1-a876-e8ddc11dcf21/field_alternative_title?resourceVersion=id%3A48"
//      },
//      "self": {
//        "href": "http://islandora-idc.traefik.me/jsonapi/node/islandora_object/815a4c04-0be5-44f1-a876-e8ddc11dcf21/relationships/field_alternative_title?resourceVersion=id%3A48"
//      }
//    }
//  }
type JsonApiLanguageValue struct {
	JsonApiData
	Meta struct {
		Value string
	}
}

// Answers the language code of the value string by resolving the Language Taxonomy entity identified in the
// JsonApiLanguageValue
func (lv JsonApiLanguageValue) LangCode(t *testing.T) string {
	jsonApiLang := JsonApiLanguage{}
	lv.Resolve(t, &jsonApiLang)
	return jsonApiLang.JsonApiData[0].JsonApiAttributes.LanguageCode
}

// Answers the value of the string, the language of which is provided by langCode(...)
func (lv JsonApiLanguageValue) Value() string {
	return lv.Meta.Value
}

// Represents the results of a JSONAPI query for a single Corporate Body Term
type JsonApiCorporateBody struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Name        string
			Description struct {
				Value     string
				Format    string
				Processed string
			}
			Authority []struct {
				Uri    string
				Title  string
				Source string
			} `json:"field_authority_link"`
			PrimaryName     string   `json:"field_primary_name"`
			SubordinateName []string `json:"field_subordinate_name"`
			Location        []string `json:"field_location_of_meeting"`
			NumberOrSection []string `json:"field_num_of_section_or_meet"`
			DateOfMeeting   []string `json:"field_date_of_meeting_or_treaty"`
			AltName         []string `json:"field_corporate_body_alt_name"`
			Date            []string `json:"field_date"`
		} `json:"attributes"`
		JsonApiRelationships struct {
			Relationships struct {
				Data []struct {
					JsonApiData
					Meta map[string]string
				}
			} `json:"field_relationships"`
		} `json:"relationships"`
	} `json:"data"`
}

type JsonApiIslandoraModel struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Name        string
			Description struct {
				Value     string
				Format    string
				Processed string
			}
			ExternalUri struct {
				Uri   string
				Title string
			} `json:"field_external_uri"`
		} `json:"attributes"`
	} `json:"data"`
}

type JsonApiIslandoraDisplay struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Name        string
			Description struct {
				Value     string
				Format    string
				Processed string
			}
			ExternalUri struct {
				Uri   string
				Title string
			} `json:"field_external_uri"`
		} `json:"attributes"`
	} `json:"data"`
}

type RelData struct {
	JsonApiData
	Meta map[string]interface{}
}

type RelContributor struct {
	Data []RelData
}

var ErrConversion = errors.New("cannot convert type")
var ErrMissing = errors.New("missing field from meta")

func (rd RelData) MetaString(field string) (string, error) {
	if value, exists := rd.Meta[field]; exists {
		if strValue, ok := value.(string); ok {
			return strValue, nil
		} else {
			return "", fmt.Errorf("%w: %v to string", ErrConversion, value)
		}
	}

	return "", fmt.Errorf("%w: %s", ErrMissing, field)
}

func (rd RelData) MetaInt(field string) (int, error) {
	if value, exists := rd.Meta[field]; exists {
		if intVal, ok := value.(int); ok {
			return intVal, nil
		} else {
			return -1, fmt.Errorf("%w: %v to int", ErrConversion, value)
		}
	}

	return -1, fmt.Errorf("%w: %s", ErrMissing, field)
}

// https://islandora-idc.traefik.me/jsonapi/media/image?filter[id]=090690a5-4db5-4d72-a94e-3b26a90b516b
type JsonApiImageMedia struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			JsonApiMediaAttributes
			JsonApiImageMediaAttributes
		} `json:"attributes"`
		JsonApiRelationships struct {
			JsonApiMediaRelationships
			File struct {
				Data RelData
			} `json:"field_media_image"`
		} `json:"relationships"`
	} `json:"data"`
}

type JsonApiMediaAttributes struct {
	FileSize         int    `json:"field_file_size"`
	MimeType         string `json:"field_mime_type"`
	OriginalName     string `json:"field_original_name"`
	Name             string
	RestrictedAccess bool `json:"field_restricted_access"`
}

type JsonApiMediaRelationships struct {
	AccessTerms struct {
		Data []JsonApiData
	} `json:"field_access_terms"`
	MediaUse struct {
		Data []JsonApiData
	} `json:"field_media_use"`
	MediaOf struct {
		Data JsonApiData
	} `json:"field_media_of"`
}

type JsonApiImageMediaAttributes struct {
	Height int `json:"field_height"`
	Width  int `json:"field_width"`
}

type JsonApiDocumentMedia struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			JsonApiMediaAttributes
		} `json:"attributes"`
		JsonApiRelationships struct {
			JsonApiMediaRelationships
			File struct {
				Data RelData
			} `json:"field_media_document"`
		} `json:"relationships"`
	} `json:"data"`
}

type JsonApiAudioMedia struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			JsonApiMediaAttributes
		} `json:"attributes"`
		JsonApiRelationships struct {
			JsonApiMediaRelationships
			File struct {
				Data RelData
			} `json:"field_media_audio_file"`
		} `json:"relationships"`
	} `json:"data"`
}

type JsonApiExtractedTextMedia struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			JsonApiMediaAttributes
			JsonApiExtractedTextMediaAttributes
		} `json:"attributes"`
		JsonApiRelationships struct {
			JsonApiMediaRelationships
			File struct {
				Data RelData
			} `json:"field_media_file"`
		} `json:"relationships"`
	} `json:"data"`
}

type JsonApiExtractedTextMediaAttributes struct {
	EditedText struct {
		Value     string
		Format    string
		Processed string
	} `json:"field_edited_text"`
}

type JsonApiGenericFileMedia struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			JsonApiMediaAttributes
		} `json:"attributes"`
		JsonApiRelationships struct {
			JsonApiMediaRelationships
			File struct {
				Data RelData
			} `json:"field_media_file"`
		} `json:"relationships"`
	} `json:"data"`
}

type JsonApiRemoteVideoMedia struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Name             string
			EmbedUrl         string `json:"field_media_oembed_video"`
			RestrictedAccess bool   `json:"field_restricted_access"`
		} `json:"attributes"`
		JsonApiRelationships struct {
			JsonApiMediaRelationships
		} `json:"relationships"`
	} `json:"data"`
}

type JsonApiVideoMedia struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			JsonApiMediaAttributes
		} `json:"attributes"`
		JsonApiRelationships struct {
			JsonApiMediaRelationships
			File struct {
				Data RelData
			} `json:"field_media_video_file"`
		} `json:"relationships"`
	} `json:"data"`
}

type JsonApiFile struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Filename string
			Uri      struct {
				Url   string
				Value string
			}
			MimeType    string `json:"filemime"`
			FileSize    int
			CreatedDate string `json:"created"`
			ChangedDate string `json:"changed"`
		} `json:"attributes"`
	} `json:"data"`
}

type JsonApiMediaUse struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			Name        string
			Description struct {
				Value     string
				Format    string
				Processed string
			}
			ExternalUri struct {
				Uri   string
				Title string
			} `json:"field_external_uri"`
		} `json:"attributes"`
		JsonApiRelationships struct {
		} `json:"relationships"`
	} `json:"data"`
}

type JsonApiFitsMedia struct {
	JsonApiData []struct {
		Type              jsonapi.DrupalType
		Id                string
		JsonApiAttributes struct {
			JsonApiMediaAttributes
		} `json:"attributes"`
		JsonApiRelationships struct {
			JsonApiMediaRelationships
			File struct {
				Data RelData
			} `json:"field_media_file"`
		} `json:"relationships"`
	} `json:"data"`
}
