package com.demo.pdfbox;

import org.apache.pdfbox.contentstream.PDFGraphicsStreamEngine;
import org.apache.pdfbox.cos.COSName;
import org.apache.pdfbox.pdmodel.graphics.image.PDImage;
import org.apache.pdfbox.pdmodel.interactive.annotation.PDAnnotationWidget;
import org.apache.pdfbox.pdmodel.interactive.annotation.PDAppearanceStream;

import java.awt.geom.Point2D;
import java.io.IOException;

public class WidgetImageChecker extends PDFGraphicsStreamEngine {
    WidgetImageChecker(PDAnnotationWidget widget) {
        super(widget.getPage());
        this.widget = widget;
    }

    boolean hasImages() throws IOException {
        count = 0;
        PDAppearanceStream normalAppearance = widget.getNormalAppearanceStream();
        processChildStream(normalAppearance, widget.getPage());
        return count != 0;
    }

    @Override
    public void drawImage(PDImage pdImage) throws IOException {
        count++;
    }

    @Override
    public void appendRectangle(Point2D p0, Point2D p1, Point2D p2, Point2D p3) throws IOException {
    }

    @Override
    public void clip(int windingRule) throws IOException {
    }

    @Override
    public void moveTo(float x, float y) throws IOException {
    }

    @Override
    public void lineTo(float x, float y) throws IOException {
    }

    @Override
    public void curveTo(float x1, float y1, float x2, float y2, float x3, float y3) throws IOException {
    }

    @Override
    public Point2D getCurrentPoint() throws IOException {
        return null;
    }

    @Override
    public void closePath() throws IOException {
    }

    @Override
    public void endPath() throws IOException {
    }

    @Override
    public void strokePath() throws IOException {
    }

    @Override
    public void fillPath(int windingRule) throws IOException {
    }

    @Override
    public void fillAndStrokePath(int windingRule) throws IOException {
    }

    @Override
    public void shadingFill(COSName shadingName) throws IOException {
    }

    final PDAnnotationWidget widget;
    int count = 0;
}
