"use strict";

const express = require("express");
const app = express();

const port = 8080;
const host = "0.0.0.0";

// Tiện học nodejs luôn =))
app.get("/", (req, res) => {
  res.send("Hello World from IBM Cloud Essentials!");
});

app.listen(port, host);
console.log(`Running on http://${host}:${port}`);

const verifyPDF = require("@ninja-labs/verify-pdf");
const fs = require("fs");
const signedPdfBuffer = fs.readFileSync("signed.pdf");

const verifyResult = verifyPDF(signedPdfBuffer);
console.log(verifyResult.signatures[0].meta.certs[0].issuedBy);
