package shield

import (
	"github.com/reiver/go-porterstemmer"
	"regexp"
	"strconv"
)

const useStemmer bool = true

type enTokenizer struct {
}

func NewEnglishTokenizer() Tokenizer {
	return &enTokenizer{}
}

func (t *enTokenizer) Tokenize(text string) (words map[string]int64) {
	words = make(map[string]int64)
	pv_str := ""
	for _, w := range splitTokenRx.Split(text, -1) {
		if len(w) > 2 {

			//*** sitnan ***** adding stemmer/
			//bigram
			if _, err := strconv.Atoi(w); err != nil { // ignore numeric values
				if !isStopWord(w, stopwords) { // ignore stop words
					stem := porterstemmer.StemString(w)
					if len(stem) > 2 {

						// trying ignoring common words during test ( do not do this while training)
						//if !isStopWord(stem, common_words) {

						words[stem]++
						if "" != pv_str {
							bg_word := pv_str + "_" + stem
							//bigram_words[bg_word]++
							words[bg_word]++
						}
						pv_str = stem
						//}

					}
				}
			}
		}
	}

	return
}

func isStopWord(s string, list []string) bool {
	for _, b := range list {
		if b == s {
			return true
		}
	}
	return false
}

// Spamassassin stoplist
//
// http://wiki.apache.org/spamassassin/BayesStopList
//
//var splitTokenRx = regexp.MustCompile(`[^\w]+|news|htm|html|article|2015|2016|who|what|where|which|how|why|when|his|her|able|all|already|and|any|are|because|both|can|come|each|email|even|few|first|for|from|give|has|have|http|information|into|it's|just|know|like|long|look|made|mail|mailing|mailto|make|many|more|most|much|need|not|now|number|off|one|only|out|own|people|place|right|same|see|such|that|the|this|through|time|using|web|where|why|with|without|work|world|year|years|you|you're|your`)
var splitTokenRx = regexp.MustCompile(`[^\w]+|_`)
var stopwords = []string{"http", "https", "news", "hindi", "meaning", "photo", "video", "slideshow", "htm", "html", "article", "able", "about", "above", "abst", "accordance", "according", "accordingly", "across", "act", "actually", "added", "adj", "affected", "affecting", "affects", "after", "afterwards", "again", "against", "ah", "all", "almost", "alone", "along", "already", "also", "although", "always", "am", "among", "amongst", "an", "and", "announce", "another", "any", "anybody", "anyhow", "anymore", "anyone", "anything", "anyway", "anyways", "anywhere", "apparently", "approximately", "are", "aren", "arent", "arise", "around", "as", "aside", "ask", "asking", "at", "auth", "available", "away", "awfully", "b", "back", "be", "became", "because", "become", "becomes", "becoming", "been", "before", "beforehand", "begin", "beginning", "beginnings", "begins", "behind", "being", "believe", "below", "beside", "besides", "better", "between", "beyond", "biol", "both", "brief", "briefly", "but", "by", "c", "ca", "came", "can", "cannot", "cant", "cause", "causes", "certain", "certainly", "co", "com", "come", "comes", "contain", "containing", "contains", "could", "couldnt", "d", "date", "did", "didnt", "different", "do", "does", "doesnt", "doing", "done", "dont", "down", "downwards", "due", "during", "e", "each", "ed", "edu", "effect", "eg", "eight", "eighty", "either", "else", "elsewhere", "end", "ending", "enough", "especially", "et", "et-al", "etc", "even", "ever", "every", "everybody", "everyone", "everything", "everywhere", "ex", "except", "f", "far", "few", "ff", "fifth", "first", "five", "fix", "followed", "following", "follows", "for", "former", "formerly", "forth", "found", "four", "from", "further", "furthermore", "g", "gave", "get", "gets", "getting", "give", "given", "gives", "giving", "go", "goes", "gone", "got", "gotten", "h", "had", "happens", "hardly", "has", "hasnt", "have", "havent", "having", "he", "hed", "hence", "her", "here", "hereafter", "hereby", "herein", "heres", "hereupon", "hers", "herself", "hes", "hi", "hid", "him", "himself", "his", "hither", "home", "how", "howbeit", "however", "hundred", "i", "id", "ie", "if", "ill", "im", "immediate", "immediately", "importance", "important", "in", "inc", "indeed", "index", "information", "instead", "into", "invention", "inward", "is", "isnt", "it", "itd", "itll", "its", "itself", "ive", "j", "just", "k", "keep", "keeps", "kept", "kg", "km", "know", "known", "knows", "l", "largely", "last", "lately", "later", "latter", "latterly", "least", "less", "lest", "let", "lets", "like", "liked", "likely", "line", "little", "ll", "look", "looking", "looks", "ltd", "m", "made", "mainly", "make", "makes", "many", "may", "maybe", "me", "mean", "means", "meantime", "meanwhile", "merely", "mg", "might", "million", "miss", "ml", "more", "moreover", "most", "mostly", "mr", "mrs", "much", "mug", "must", "my", "myself", "n", "na", "name", "namely", "nay", "nd", "near", "nearly", "necessarily", "necessary", "need", "needs", "neither", "never", "nevertheless", "new", "next", "nine", "ninety", "no", "nobody", "non", "none", "nonetheless", "noone", "nor", "normally", "nos", "not", "noted", "nothing", "now", "nowhere", "o", "obtain", "obtained", "obviously", "of", "off", "often", "oh", "ok", "okay", "old", "omitted", "on", "once", "one", "ones", "only", "onto", "or", "ord", "other", "others", "otherwise", "ought", "our", "ours", "ourselves", "out", "outside", "over", "overall", "owing", "own", "p", "page", "pages", "part", "particular", "particularly", "past", "per", "perhaps", "placed", "please", "plus", "poorly", "possible", "possibly", "potentially", "pp", "predominantly", "present", "previously", "primarily", "probably", "promptly", "proud", "provides", "put", "q", "que", "quickly", "quite", "qv", "r", "ran", "rather", "rd", "re", "readily", "really", "recent", "recently", "ref", "refs", "regarding", "regardless", "regards", "related", "relatively", "research", "respectively", "resulted", "resulting", "results", "right", "run", "s", "said", "same", "saw", "say", "saying", "says", "sec", "section", "see", "seeing", "seem", "seemed", "seeming", "seems", "seen", "self", "selves", "sent", "seven", "several", "shall", "she", "shed", "shell", "shes", "should", "shouldnt", "show", "showed", "shown", "showns", "shows", "significant", "significantly", "similar", "similarly", "since", "six", "slightly", "so", "some", "somebody", "somehow", "someone", "somethan", "something", "sometime", "sometimes", "somewhat", "somewhere", "soon", "sorry", "specifically", "specified", "specify", "specifying", "still", "stop", "strongly", "sub", "substantially", "successfully", "such", "sufficiently", "suggest", "sup", "sure", "t", "take", "taken", "taking", "tell", "tends", "th", "than", "thank", "thanks", "thanx", "that", "thatll", "thats", "thatve", "the", "their", "theirs", "them", "themselves", "then", "thence", "there", "thereafter", "thereby", "thered", "therefore", "therein", "therell", "thereof", "therere", "theres", "thereto", "thereupon", "thereve", "these", "they", "theyd", "theyll", "theyre", "theyve", "thing", "things", "think", "this", "those", "thou", "though", "thoughh", "thousand", "throug", "through", "throughout", "thru", "thus", "til", "tip", "to", "together", "too", "took", "toward", "towards", "tried", "tries", "truly", "try", "trying", "ts", "twice", "two", "u", "un", "under", "unfortunately", "unless", "unlike", "unlikely", "until", "unto", "up", "upon", "ups", "us", "use", "used", "useful", "usefully", "usefulness", "uses", "using", "usually", "v", "value", "various", "ve", "very", "via", "viz", "vol", "vols", "vs", "w", "want", "wants", "was", "wasnt", "way", "we", "wed", "welcome", "well", "went", "were", "werent", "weve", "what", "whatever", "whatll", "whats", "when", "whence", "whenever", "where", "whereafter", "whereas", "whereby", "wherein", "wheres", "whereupon", "wherever", "whether", "which", "while", "whim", "whither", "who", "whod", "whoever", "whole", "wholl", "whom", "whomever", "whos", "whose", "why", "widely", "will", "willing", "wish", "with", "within", "without", "wont", "words", "world", "would", "wouldnt", "www", "x", "y", "yes", "yet", "you", "youd", "youll", "your", "youre", "yours", "yourself", "yourselves", "youve", "z", "zero"}

var common_words = []string{"actress", "andhra", "andhra_pradesh", "arrest", "articleshow", "attack", "bangalor", "beauti", "bihar", "bjp", "bollywood", "case", "cast", "celebr", "chennai", "citi", "court", "dai", "delhi", "ent", "entertain", "facebook", "fbo", "featur", "galleri", "girl", "gossip", "health", "hmu", "hot", "india", "indian", "int", "intern", "kapoor", "karnataka", "kerala", "khan", "kill", "leader", "life", "man", "modi", "movi", "mumbai", "nadu", "nation", "nor", "offic", "omc", "pakistan", "pari", "peopl", "pho", "photo", "photogalleri", "photogalleri_entertain", "pic", "polic", "polit", "polit_nation", "pradesh", "refer", "relationship", "review", "road", "sex", "singh", "slideshow", "special", "state", "stori", "student", "style", "tamil", "tamilnadu", "tamil_nadu", "top", "topic", "woman", "women", "year"}
