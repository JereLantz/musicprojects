//create quick buttons
function quickButtonCreation(){
    const elementCount = document.getElementsByClassName("progSelectors")
    const quickButtons = document.getElementsByClassName("quickButtons")

    switch(elementCount.length){
        case 4:
            let parent1625 = document.getElementById("quickButtons")
            let createdBtn1625 = document.createElement("button")

            if(document.getElementById("set251")){
                document.getElementById("set251").remove()
            }

            if(!(document.getElementById("set1625"))){
                createdBtn1625.id = "set1625"
                createdBtn1625.textContent = "Set 1625"
                createdBtn1625.className = "quickButtons"

                parent1625.appendChild(createdBtn1625)

                document.getElementById("set1625").addEventListener("click", set1625)
            }
            break
        case 3:
            let parent251 = document.getElementById("quickButtons")
            let createdBtn251 = document.createElement("button")

            if(document.getElementById("set1625")){
                document.getElementById("set1625").remove()
            }

            if(!(document.getElementById("set251"))){
                createdBtn251.id = "set251"
                createdBtn251.textContent = "Set 251"
                createdBtn251.className = "quickButtons"

                parent251.appendChild(createdBtn251)
                document.getElementById("set251").addEventListener("click", set251)
            }
            break;

        default:
            for(button of quickButtons){
                button.remove()
            }

            break
    }
}

// button to set che chord progress to 251
function set251(){
    document.getElementById("progSelect1").value = 2
    document.getElementById("progSelect2").value = 5
    document.getElementById("progSelect3").value = 1
}

// Button to set the chord progress to 1625
function set1625(){
    document.getElementById("progSelect1").value = 1
    document.getElementById("progSelect2").value = 6
    document.getElementById("progSelect3").value = 2
    document.getElementById("progSelect4").value = 5
}

// Get the chord prog from the dropdown menus when the button is pressed
function getProgression(){
    // Get the key of the wanted chord progress
    const theKey = document.getElementById("keySelect").value
    const theMode = document.getElementById("modeSelect").value
    const progElements = document.getElementsByClassName("progSelectors");
    // array for the chords in the progression
    let chords = ""
    
    const progressionInput = [];

    for(let i = 0; i < progElements.length; i++){
        let valueFromElement = document.getElementById(`progSelect${i+1}`).value

        progressionInput.push(valueFromElement);
    }

    for(const chord of progressionInput){
        switch(theKey){
            // Note: The keys F#, C#, Cb are not implemented at the moment
            // Keys with sharps in them
            case "G":
            case "D":
            case "B":
            case "E":
            case "A":
                chords += sharpKeys(theKey, chord, theMode);
                chords += chordQuality(chord, theMode);
                chords += " ";
                break;

            // Keys with flats in them
            case "C": // C Major doesn't have sharps or flats, but the modes of C have flats
            case "F":
            case "Bb":
            case "Eb":
            case "Ab":
            case "Db":
            case "Gb":
                chords += flatKeys(theKey, chord, theMode);
                chords += chordQuality(chord, theMode);
                chords += " ";
                break;
        }
    
    }

    document.getElementById("displayProgText").innerText = `In the key of "${theKey}" the progression ${progressionInput} is:`
    document.getElementById("displayProgression").innerText = chords
    displayOutput()
}

/**
 * Creates/removes the progress select elements based on user how many the user wants to be visible
 */
document.getElementById("progLenBtn").addEventListener("click", setProgLen);
function setProgLen(){
    const currentElements = document.getElementsByClassName("progSelectors");
    const wantedLen = document.getElementById("lenInputNbr").value;

    if(wantedLen > currentElements.length){
        for(let i = currentElements.length; i < wantedLen; i++){
            let parent = document.getElementById("progSelectors");
            let selectToCreate = document.createElement("select");

            selectToCreate.id = `progSelect${i+1}`;
            selectToCreate.className = "progSelectors"
            parent.appendChild(selectToCreate);

            for(let i = 0; i < 7; i++){
                let option = document.createElement("option");

                option.value = i + 1;
                option.text = i + 1;
                
                selectToCreate.appendChild(option);
            }
        }
    }
    else if(wantedLen < currentElements.length){
        for(let i = currentElements.length; i > wantedLen; i--){
            currentElements[i-1].remove();
        }
    }
    quickButtonCreation();
    
    if(!(document.getElementById("progBtn"))){
        createGetProgBtn()
    }
}

/**
 * Creates the "get progress" button
 */
function createGetProgBtn(){
    let parent = document.getElementById("getProgBtn")
    let getProgBtn = document.createElement("button")

    getProgBtn.id = "progBtn"
    getProgBtn.textContent = "Get progression"

    parent.appendChild(getProgBtn)
    document.getElementById("progBtn").addEventListener("click", getProgression)
}

// Makes the elements in displayOutput visible
function displayOutput(){
    for(const element of document.getElementsByClassName("displayOutput")){
        element.style.visibility = "visible"
    }
}

// Keys with sharps in them
function sharpKeys(theKey, chord, mode){
    const possibleNotes = ["C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"];
    const noteDistances = [2, 2, 1, 2, 2, 2, 1];
    let noteFromArr = 0;
    let i = 0;

    for(const note of possibleNotes){
        if(note == theKey){
            noteFromArr = i;
            break;
        }
        i++;
    }

    for(i = Number(mode); i < Number(mode) + Number(chord)-1; i++){
        noteFromArr += noteDistances[i % 7];
    }

    return possibleNotes[noteFromArr%12];
}

// Keys with flats in them
function flatKeys(theKey, chord, mode){
    const possibleNotes = ["C", "Db", "D", "Eb", "E", "F", "Gb", "G", "Ab", "A", "Bb", "B"];
    const noteDistances = [2, 2, 1, 2, 2, 2, 1];
    let noteFromArr = 0;
    let i = 0;

    for(const note of possibleNotes){
        if(note == theKey){
            noteFromArr = i;
            break;
        }
        i++;
    }

    for(i = Number(mode); i < Number(mode) + Number(chord)-1; i++){
        noteFromArr += noteDistances[i % 7];
    }
    return possibleNotes[noteFromArr%12];
}

// Returns the quality of the chord based on the key
function chordQuality(chordNbr, mode){
    let chordCalcToScale = ((Number(chordNbr) - 1 + Number(mode)) %7 + 1);

    if(Number(mode) == 3){
        return lydianQuality(chordNbr)
    }

    switch(Number(chordCalcToScale)){
        case 1:
        case 4:
            return "Maj7";
        
        case 2:
        case 3:
        case 6:
            return "m7";
        
        case 5:
            return "7";
            
        case 7:
            return "m7b5";
    }
}

// Lydian mode is special, so it has to be done separately
function lydianQuality(chordNbr){
    switch(Number(chordNbr)){
        case 1:
            return "Maj7";
        case 2:
            return "7";
        case 3:
        case 6:
        case 7:
            return "m7";
        case 4:
            return "m7b5";
        case 5:
            return "6";
    }
}