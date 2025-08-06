//@ts-check

function setSelectorsNatural(){
    const selectors = document.querySelectorAll("select")

    for(let select of selectors){
        select.value = "natural"
    }
}

function setSelectorsSharp(){
    const selectors = document.querySelectorAll("select")

    for(let select of selectors){
        select.value = "sharp"
    }
}

function setSelectorsFlat(){
    const selectors = document.querySelectorAll("select")

    for(let select of selectors){
        select.value = "flat"
    }
}

/** adds a chord selector to the chord progress page
 */
function addChordSelect(){
    const parent = document.getElementById("chordSelectors")

    const selects = parent?.querySelectorAll("select")
    if (!selects) return

    if (selects?.length > 8){
        return
    }

    const newSel = document.createElement("select")
    newSel.name = "chordSelect"
    for (let i = 1; i < 8; i++){
        const newOption = document.createElement("option")
        newOption.value = String(i)
        newOption.innerText = String(i)

        newSel.appendChild(newOption)
    }

    parent?.appendChild(newSel)
}

/** removes a chord selector on the chord progress page
 */
function removeChordSelect(){
    const parent = document.getElementById("chordSelectors")
    if (!parent) return

    const selects = parent?.querySelectorAll("select")
    if(!selects) return

    if(selects.length > 1) {
        selects[selects.length - 1].remove()
    }
}
