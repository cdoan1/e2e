# This python script iterates the list of junit results, and creates an aggregate csv of the results
# into the ./data/results.csv file
# TODO: If s3 credentials around found. we'll push the results.csv to s3
#
import xml.etree.ElementTree as ET
import csv
import glob, os

# change dir to data
os.chdir("./data")

csv_filename = open("results.csv", 'w')
csv_writer = csv.writer(csv_filename)

head = ['date','login','overview','toplogy','cluster','application','policy']
csv_writer.writerow(head)

for file in glob.glob("*.xml"):
    filename = file.split(".")[0]
    tree = ET.parse(filename + ".results.xml")
    root = tree.getroot()

    row = []
    row.append(filename)

    for member in root.findall('testcase'):
        data = member.attrib
        row.append(data['time'])

    csv_writer.writerow(row)

csv_filename.close()
