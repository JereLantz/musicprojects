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
