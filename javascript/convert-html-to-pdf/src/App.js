import "./App.css";
import { useRef } from "react";
import ReportTemplate from "./ReportTemplate";
import html2PDF from "jspdf-html2canvas";

function App() {
  const reportTemplateRef = useRef(null);

  const handleGeneratePdf = () => {
    html2PDF(reportTemplateRef.current, {
      jsPDF: {
        format: "a4",
        orientation: "p",
      },
      success: function (pdf) {
        window.open(pdf.output("bloburl"));
      },
    });
  };

  return (
    <div>
      <button className="button" onClick={handleGeneratePdf}>
        Generate PDF
      </button>
      <div ref={reportTemplateRef}>
        <ReportTemplate />
      </div>
    </div>
  );
}

export default App;
