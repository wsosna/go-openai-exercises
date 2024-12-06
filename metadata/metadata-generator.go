package metadata

import (
	"encoding/json"
	"github.com/openai/openai-go"
	client2 "go-openai-exercises/client"
	my_ai "go-openai-exercises/my-ai"
	"go-openai-exercises/utils"
	"log"
	"strings"
)

const AllFormatsDir = "tmp/pliki_z_fabryki"

func Run() {
	log.Println("Running metadata exercise")
	facts := utils.ReadFilesToPrompt(AllFormatsDir + "/facts")
	reports := utils.ReadFilesToPrompt(AllFormatsDir+"/reports", "report")

	client := my_ai.NewOpenAiWrapper(openai.ChatModelGPT4o)
	system := `
<system>
You are a helpful assistant that assigns keywords to files based on their content.
</system>
`

	user := `
<message>
Generate a map of the following files. Map should contain the name of a mentioned person and keywords from the file about that person.
Map should be generated in Polish, same as the content of the files.
<files>
` + facts + `
</files>
</message>
`

	mapOfFacts := client.AskMyAI(system, user)
	log.Println(mapOfFacts)

	system = `
<system>
You are a helpful assistant that assigns keywords to files based on their content.
</system>
`

	user = `
<message>
Generate tags for the following reports. Provide your response in JSON format as shown in the example below. 
Ensure that your response contains only the JSON. You should use your own knowledge and the map of facts to generate the tags.
<map_of_facts>
` + mapOfFacts + `
</map_of_facts>
<reports>
` + reports + `
</reports>

Respond only with the result in JSON format, as shown in the example below:
<example>
{
"nazwa-pliku-01.txt":"lista, słów, kluczowych 1",
"nazwa-pliku-02.txt":"lista, słów, kluczowych 2",
"nazwa-pliku-03.txt":"lista, słów, kluczowych 3",
"nazwa-pliku-NN.txt":"lista, słów, kluczowych N"
}
</example>
</message>
`

	resp := client.AskMyAI(system, user)
	log.Println("Received response from AI")
	resp = strings.ReplaceAll(resp, "```json", "")
	resp = strings.ReplaceAll(resp, "```", "")
	log.Println(resp)

	centrala := client2.Centrala{}

	log.Println("Sending solution to centrala")
	r := unmarshal(resp)
	resp = centrala.SendSolution("dokumenty", r)

	log.Println("Received response from centrala")
	log.Println(resp)

	askAiMock(system, user, *client)

}

func unmarshal(jsonString string) map[string]string {
	var reports map[string]string
	err := json.Unmarshal([]byte(jsonString), &reports)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	return reports
}

func askAiMock(s1 string, s2 string, openAiWrapper my_ai.OpenAiWrapper) {
	//log.Println(s1 + s2 + openAiWrapper.Model)
}

func Cheat() {
	resp := `
{
    "2024-11-12_report-00-sektor_C4.txt": "Aleksander Ragowski, nauczyciel języka angielskiego, Szkoła Podstawowa nr 9 w Grudziądzu, kreatywne metody nauczania, zaangażowanie społeczne, krytyk reżimu robotów, tajne spotkania, zagrożenia edukacyjne, programowanie w Java, poszukiwany, aktywista ruchu oporu",
    "2024-11-12_report-01-sektor_A1.txt": "lokalna zwierzyna leśna, fałszywy alarm",
    "2024-11-12_report-02-sektor_A3.txt": "brak wykrycia, nocny patrol",
    "2024-11-12_report-03-sektor_A3.txt": "czujniki aktywne, monitorowanie, brak wykrycia",
    "2024-11-12_report-04-sektor_B2.txt": "patrol, bezpieczeństwo, brak anomalii",
    "2024-11-12_report-05-sektor_C1.txt": "monitorowanie, czujniki, brak aktywności",
    "2024-11-12_report-06-sektor_C2.txt": "spokój, brak wykrycia, patrol",
    "2024-11-12_report-07-sektor_C4.txt": "Barbara Zawadzka, frontend development, branża IT, ruch oporu, JavaScript, Python, AI Devs, sztuczna inteligencja, Związek z Aleksandrem Ragowskim, walka wręcz, ultradźwiękowy sygnał, tajne technologie, sektor C4",
    "2024-11-12_report-08-sektor_A1.txt": "cisza, brak ruchu, obserwacja",
    "2024-11-12_report-09-sektor_C2.txt": "brak anomalii, zakończenie patrolu"
}
`

	centrala := client2.Centrala{}
	r := centrala.SendSolution("dokumenty", unmarshal(resp))
	log.Println(r)
}
