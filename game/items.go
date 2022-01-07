package game

type Item struct {
	Entity
	Sequence    string
	Description string
}

func Protein(p Pos) *Item {
	return &Item{Entity{p, "Protein", 'P'}, "Methionine-Alanine-Stop", ""}

}
func RNAPol(p Pos) *Item {
	return &Item{Entity{p, "RNA Polymerase", 'E'}, "", ""}

}
func Instructions(p Pos, level *Level) *Item {
	return &Item{Entity{p, "Instructions", 'I'}, "", "1. Welcome " + level.Player.Name + " you have a gene from the DNA that you would like to turn into a protein for use. The sequence is ATCCGAACT. The first thing is to transcribe the gene into mRNA "}
}

func FirstDirection(p Pos) *Item {
	return &Item{Entity{p, "First Directions", 'Q'}, "", "2. First take gene sequence and find the RNA polymerase and 3 transcription factors to continue"}
}

func GeneSequence(p Pos) *Item {
	return &Item{Entity{p, "Gene Sequence", 'S'}, "ATCCGAACT", ""}
}

func TranscriptionFactor(p Pos) *Item {
	return &Item{Entity{p, "Transcription Factor", 'T'}, "", ""}
}

func InterInfo(p Pos) *Item {
	return &Item{Entity{p, "Intermediate Info", 'N'}, "UAGGCUUGA", "Congratulations you have completed the first step of protein synthesis, transcription. With the three transcription factors and the enyme RNA Polymerase, you transcribed your DNA into UAGGCUUGA. The complementary base for Adenine(A) when transcribing into RNA is Uracil(U). Now go to the cytoplasm by going down the stairs"}
}

func TranslationInfo(p Pos) *Item {
	return &Item{Entity{p, "Translation Info", 'R'}, "", "3. You are now in the ribosome, where translation of the mRNA takes place. You have with you the mRNA sequence UAGGCUUGA. The sequence will be translated by reading the sequence 3 bases at a time and creating an amino acid based on the codon table (Press C). By chaining amino acids, you get a protein. Try figuring out the protein chain"}
}

func FinalNote(p Pos) *Item {
	return &Item{Entity{p, "Final Note", 'F'}, "", "So in this game you learned that protein synthesis begins with a gene sequence (what you got at the beginning). Using transcription factors and RNA polymerase, you transcribe the gene into an mRNA sequence in the nucleus. Then the mRNA sequence gets translated into an amino acid chain, creating a protein, in the ribosomes. There is more to learn than that, but that was a brief overview. Thank you so much for playing my game!"}
}
