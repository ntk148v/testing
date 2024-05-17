import "./App.css";
import { useRef } from "react";
import ReportTemplate from "./ReportTemplate";
import html2PDF from "jspdf-html2canvas";
import jsPDF from "jspdf";

function App() {
  const reportTemplateRef = useRef(null);

  const handleGeneratePdfWithjsPDF = () => {
    const contentWidth = reportTemplateRef.current.scrollWidth;
    const contentHeight = reportTemplateRef.current.scrollHeight;
    const orientation = contentWidth > contentHeight ? "l" : "p";
    const doc = new jsPDF({
      orientation: orientation,
      format: "a4",
    });

    const ratio =
      contentWidth > contentHeight
        ? doc.internal.pageSize.getWidth() / contentWidth
        : doc.internal.pageSize.getHeight() / contentHeight;

    doc.html(reportTemplateRef.current, {
      async callback(pdf) {
        window.open(pdf.output("bloburl"));
      },
      html2canvas: {
        scale: ratio,
        //600 is the width of a4 page. 'a4': [595.28, 841.89]
      },
    });
  };

  const handleGeneratePdfWithHtml2PDF = () => {
    const contentWidth = reportTemplateRef.current.scrollWidth;
    const contentHeight = reportTemplateRef.current.scrollHeight;
    const orientation = contentWidth > contentHeight ? "l" : "p";
    html2PDF(reportTemplateRef.current, {
      jsPDF: {
        format: "a4",
        orientation: orientation,
      },
      success: function (pdf) {
        window.open(pdf.output("bloburl"));
      },
    });
  };

  return (
    <div>
      <button className="button" onClick={handleGeneratePdfWithjsPDF}>
        Generate PDF 1
      </button>
      <button className="button" onClick={handleGeneratePdfWithHtml2PDF}>
        Generate PDF 2
      </button>
      <div ref={reportTemplateRef}>
        <ReportTemplate />
      </div>
    </div>
  );
}

export default App;
