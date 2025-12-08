# ordi

Terminal-Anwendung zur Dateiverwaltung und Bildersortierung 

## Installation

### Voraussetzungen
- Go 1.16 oder höher

### Build
```bash
# Repository klonen
git clone https://github.com/killerpanda03/mindfocus.git
cd ordi

# Dependencies installieren
go mod download

# Program lokal ausführen
go run main.go

# Anwendung bauen
go build -o ordi.exe
```

### Funktionen

1. **Ein Verzeichnis organisieren**
   - Organisiert Dateien nach Typ in kategorisierte Ordner
   - Kategorien: Bilder, Videos, Musik, Dokumente, Archive, Sonstiges

2. **Duplikate finden**
   - Findet und verwaltet doppelte Dateien in einem Verzeichnis

### Kommende Funktion

3. **Dateien komprimieren**
   - Komprimiert verschiedene Dateitypen (Bilder, Videos, Audio, PDFs, Dokumente)
   - Benötigt externe Tools (optional):
     - **ffmpeg** - für Video- und Audio-Komprimierung
     - **ImageMagick** - für Bild-Komprimierung
     - **Ghostscript** - für PDF-Komprimierung
     - **7zip** - für Dokument- und Archiv-Komprimierung
