// Run this script from the current directory
const fs = require("fs")
const path = require("path")
const mysql = require("mysql")

const imagesPath = path.join(__dirname, "../app/public/images")

const trials = [3, 4, 5, 6]

trials.forEach((trialNum)=>{
    const trialPath = path.join(imagesPath, "trial-" + trialNum)
    fs.readdir(trialPath, (err, files)=>{
        files.forEach(file=>{
            const filePath = path.join(trialPath, file)
            console.log(fs.statSync(filePath).mtime)
        })
    })
})