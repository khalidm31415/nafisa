package elasticsarch_helper

var Index = `{
	"mappings":{
        "properties":{
            "user_id":{
                "type":"keyword"
            },
            "year_born":{
                "type":"short"
            },
            "sex":{
                "type":"keyword"
            },
            "last_education":{
                "type":"keyword"
            },
            "summary":{
                "type":"text"
            },
            "summary_dense_vector":{
                "type":"dense_vector",
                "dims": 512
            }
        }
	}
}`
