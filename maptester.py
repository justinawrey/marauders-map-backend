#!/usr/bin/python3

from PyQt5.QtWidgets import QWidget, QApplication
from PyQt5.QtGui import QPainter, QColor, QBrush, QPoint
import sys
import json

USAGE_STR = "\
usage:\n\
maptester.py mapper\n\
maptest.py (density | heatmap | raw) <data.json>\n"


class Coordinate:
    DEGREE_SIZE_PX = 7

    def __init__(self, long, lat):
        self.long = long
        self.lat = lat
        self.selected = False

    def getPixelLocation(self):
        if self.long < 0:
            x = 180 - abs(self.long)
        else:
            x = 180 + self.long

        if self.lat < 0:
            y = 90 + abs(self.lat)
        else:
            y = 90 - self.lat

        return (x, y)


class RawViewer(QWidget):

    def __init__(self, jsonData):
        super().__init__()
        self.initUI()
        self.data = jsonData

    def initUI(self):
        self.setGeometry(100,
                         100,
                         Coordinate.DEGREE_SIZE_PX * 360,
                         Coordinate.DEGREE_SIZE_PX * 180)
        self.setWindowTitle('Raw Locations')
        self.show()

    def paintEvent(self, e):
        qp = QPainter()
        qp.begin(self)
        self.drawRectangles(qp)
        qp.end()

    def drawRectangles(self, qp):
        qp.setBrush(QColor(200, 0, 0))
        for jsonObj in self.data:
            longitude = jsonObj["longitude"]
            latitude = jsonObj["latitude"]
            pixelCoords = Coordinate(longitude, latitude).getPixelLocation()
            qp.drawRect(pixelCoords[0],
                        pixelCoords[1],
                        Coordinate.DEGREE_SIZE_PX,
                        Coordinate.DEGREE_SIZE_PX)


class DensityViewer(QWidget):

    def __init__(self, jsonData):
        super().__init__()
        self.initUI()
        self.data = jsonData

    def initUI(self):
        self.setGeometry(100,
                         100,
                         Coordinate.DEGREE_SIZE_PX * 360,
                         Coordinate.DEGREE_SIZE_PX * 180)
        self.setWindowTitle('Density')
        self.show()

    def paintEvent(self, e):
        qp = QPainter()
        qp.begin(self)
        self.drawDensity(qp)
        qp.end()

    def drawDensity(self, qp):
        qp.setBrush(QColor(200, 0, 0))
        for jsonObj in self.data:
            longitude = jsonObj["longitude"]
            latitude = jsonObj["latitude"]
            radius = jsonObj["radius"]
            pixelCoords = Coordinate(longitude, latitude).getPixelLocation()
            qp.drawEllipse(QPoint(pixelCoords[0], pixelCoords[1]),
                           radius,
                           radius)


class HeatmapViewer(QWidget):

    def __init__(self, jsonData):
        super().__init__()
        self.initUI()
        self.data = jsonData

    def initUI(self):
        self.setGeometry(100,
                         100,
                         Coordinate.DEGREE_SIZE_PX * 360,
                         Coordinate.DEGREE_SIZE_PX * 180)
        self.setWindowTitle('Heatmap')
        self.show()

    def paintEvent(self, e):
        qp = QPainter()
        qp.begin(self)
        self.drawHeatmap(qp)
        qp.end()

    def drawHeatmap(self, qp):
        pass


class Mapper(QWidget):

    def __init__(self):
        super().__init__()
        self.initUI()

    def initUI(self):
        self.setGeometry(100,
                         100,
                         Coordinate.DEGREE_SIZE_PX * 360,
                         Coordinate.DEGREE_SIZE_PX * 180)
        self.setWindowTitle('Mapper')
        self.show()

def main():
    app = QApplication(sys.argv)  # setup

    args = sys.argv
    if len(args) == 2 and args[1] == "mapper":
        print("opening in mapper mode")
        _ = Mapper()

    elif len(args) == 3:
        with open(args[2]) as data_file:
            data = json.load(data_file)

        if args[1] == "density":
            print("opening in density viewer mode")
            _ = DensityViewer(data)

        elif args[1] == "heatmap":
            print("opening in heatmap viewer mode")
            _ = HeatmapViewer(data)

        elif args[1] == "raw":
            print("opening in raw viewer mode")
            _ = RawViewer(data)
    else:
        print(USAGE_STR)

    sys.exit(app.exec_())  # teardown


if __name__ == "__main__":
    main()
