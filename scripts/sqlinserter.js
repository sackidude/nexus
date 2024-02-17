// Run this script from the current directory
const fs = require("fs")
const path = require("path")
const mysql = require("mysql")

const con = mysql.createConnection({
    user: "root",
    password: JSON.parse(fs.readFileSync(__dirname + "/../secrets.json")).password,
    database: "nexus"
})

const imagesPath = path.join(__dirname, "../app/public/images")

const trials = [3, 4, 5, 6]

con.connect(err=>{
    if(err){console.error(err)}
    console.log("Connected!")
    trials.forEach((trialNum)=>{
        const trialPath = path.join(imagesPath, "trial-" + trialNum)
        const t0 = fs.statSync(path.join(trialPath, "10.jpg")).mtime
        con.query(
            `INSERT INTO Trials (trial_num, directory, start, zero_height, 1000_ml_height)
            VALUES (
                ${Number(trialNum)},
                'trial-${trialNum}',
                '${t0.toISOString().slice(0,19).replace('T', ' ')}',
                0,
                100
            )`
        )

        fs.readdir(trialPath, (err, files)=>{
            files.forEach(file=>{
                const filePath = path.join(trialPath, file)
                const timestamp = fs.statSync(filePath).mtime
                const num = Number(file.substring(0,file.length-4))

                con.query(
                    `INSERT INTO Images (filename, trial, time)
                    VALUES (
                        '${file}',
                        ${trialNum},
                        '${timestamp.toISOString().slice(0,19).replace('T', ' ')}'
                    )`
                )
            })
        })
    })
})