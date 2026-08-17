import AppKit
import Foundation

private let sampleDocuments: [(String, String)] = [
    ("Field Notes.txt", """
    A native window should feel quiet and immediate.

    This document is loaded from a real text file and displayed by NSTextView. The sidebar is an NSOutlineView, the split is NSSplitViewController, and the toolbar is NSToolbar.

    There is no HTML, JavaScript, CSS, asset navigation, or WKWebView in this editor window.
    """),
    ("Ideas.txt", """
    Ideas for the native application API

    • Treat a WebView as content, not as a window.
    • Keep toolbar item handles instead of public callback IDs.
    • Let AppKit own selection, focus, undo, scrolling, and accessibility.
    • Copy document text across the Go bridge only when it is needed.
    """),
    ("Performance.txt", """
    Performance experiment

    The important measurements are process resident memory, idle CPU, wakeups, and the absence of a WebKit content process.

    Use the commands in this example's README while the window is open and again while it is hidden.
    """),
]

private struct BenchmarkConfiguration {
    let startsHidden: Bool
    let documentBytes: Int
    let autoQuitMilliseconds: Int?
    let readyFile: String?

    init(environment: [String: String] = ProcessInfo.processInfo.environment) {
        startsHidden = environment["NATIVE_EDITOR_BENCH_HIDDEN"] == "1"
        documentBytes = max(0, Int(environment["NATIVE_EDITOR_BENCH_DOCUMENT_BYTES"] ?? "0") ?? 0)
        autoQuitMilliseconds = Int(environment["NATIVE_EDITOR_BENCH_AUTO_QUIT_MS"] ?? "")
        readyFile = environment["NATIVE_EDITOR_BENCH_READY_FILE"]
    }
}

private final class NoteFile: NSObject {
    let name: String
    let url: URL

    init(name: String, url: URL) {
        self.name = name
        self.url = url
    }
}

private final class SidebarSection: NSObject {
    let title: String
    var files: [NoteFile]

    init(title: String, files: [NoteFile]) {
        self.title = title
        self.files = files
    }
}

private final class AppDelegate: NSObject,
    NSApplicationDelegate,
    NSWindowDelegate,
    NSOutlineViewDataSource,
    NSOutlineViewDelegate,
    NSSearchFieldDelegate,
    NSTextViewDelegate,
    NSToolbarDelegate
{
    private let benchmark = BenchmarkConfiguration()
    private let splitController = NSSplitViewController()
    private let outlineView = NSOutlineView()
    private let textView = NSTextView()
    private let section = SidebarSection(title: "Text Files", files: [])
    private var allFiles: [NoteFile] = []
    private var currentFile: NoteFile?
    private var isLoading = false
    private var dirty = false
    private var window: NSWindow!
    private var statusItem: NSStatusItem!
    private var saveToolbarItem: NSToolbarItem?
    private var directory: URL!

    private let saveIdentifier = NSToolbarItem.Identifier("native-notes.save")
    private let searchIdentifier = NSToolbarItem.Identifier("native-notes.search")
    private let trackingSeparatorIdentifier = NSToolbarItem.Identifier("native-notes.sidebar-tracking-separator")

    func applicationDidFinishLaunching(_ notification: Notification) {
        NSApp.setActivationPolicy(.accessory)
        do {
            directory = try materialiseDocuments()
            try buildInterface()
            buildMenuBar()
            buildStatusItem()
            if !benchmark.startsHidden {
                showWindow(nil)
            }
            signalBenchmarkReady()
            scheduleBenchmarkExit()
            fputs("editing native text files in \(directory.path)\n", stderr)
        } catch {
            fputs("Native Notes Swift: \(error)\n", stderr)
            NSApp.terminate(nil)
        }
    }

    private func materialiseDocuments() throws -> URL {
        let base = FileManager.default.temporaryDirectory
            .appendingPathComponent("swift-native-notes-\(UUID().uuidString)", isDirectory: true)
        try FileManager.default.createDirectory(at: base, withIntermediateDirectories: true)
        for (name, contents) in sampleDocuments {
            try Data(contents.utf8).write(to: base.appendingPathComponent(name))
        }
        if benchmark.documentBytes > 0 {
            let line = Array("A native editor benchmark line with predictable UTF-8 text.\n".utf8)
            var bytes = [UInt8](repeating: 0, count: benchmark.documentBytes)
            var offset = 0
            while offset < bytes.count {
                let count = min(line.count, bytes.count - offset)
                bytes.replaceSubrange(offset..<(offset + count), with: line.prefix(count))
                offset += count
            }
            try Data(bytes).write(to: base.appendingPathComponent("00 Benchmark.txt"))
        }
        return base
    }

    private func buildInterface() throws {
        let urls = try FileManager.default.contentsOfDirectory(
            at: directory,
            includingPropertiesForKeys: nil
        ).filter { $0.pathExtension.lowercased() == "txt" }
            .sorted { $0.lastPathComponent < $1.lastPathComponent }
        allFiles = urls.map { NoteFile(name: $0.lastPathComponent, url: $0) }
        section.files = allFiles

        let sidebarController = NSViewController()
        let sidebarScroll = NSScrollView()
        sidebarScroll.drawsBackground = false
        sidebarScroll.hasVerticalScroller = true
        sidebarScroll.autohidesScrollers = true
        outlineView.headerView = nil
        outlineView.addTableColumn(NSTableColumn(identifier: NSUserInterfaceItemIdentifier("notes")))
        outlineView.outlineTableColumn = outlineView.tableColumns[0]
        outlineView.dataSource = self
        outlineView.delegate = self
        outlineView.rowSizeStyle = .medium
        outlineView.backgroundColor = .clear
        if #available(macOS 11.0, *) {
            outlineView.style = .sourceList
        }
        sidebarScroll.documentView = outlineView
        sidebarController.view = sidebarScroll

        let editorController = NSViewController()
        let editorScroll = NSScrollView()
        editorScroll.hasVerticalScroller = true
        editorScroll.hasHorizontalScroller = false
        editorScroll.autohidesScrollers = true
        editorScroll.drawsBackground = true
        editorScroll.backgroundColor = .textBackgroundColor
        textView.delegate = self
        textView.isRichText = false
        textView.isAutomaticQuoteSubstitutionEnabled = false
        textView.isAutomaticDashSubstitutionEnabled = false
        textView.allowsUndo = true
        textView.isVerticallyResizable = true
        textView.isHorizontallyResizable = false
        textView.autoresizingMask = [.width]
        textView.textContainer?.widthTracksTextView = true
        textView.textContainerInset = NSSize(width: 28, height: 28)
        textView.font = NSFont.monospacedSystemFont(ofSize: 14, weight: .regular)
        editorScroll.documentView = textView
        editorController.view = editorScroll

        let sidebarItem = NSSplitViewItem(sidebarWithViewController: sidebarController)
        sidebarItem.minimumThickness = 210
        sidebarItem.maximumThickness = 340
        sidebarItem.canCollapse = true
        let contentItem = NSSplitViewItem(viewController: editorController)
        contentItem.minimumThickness = 420
        splitController.addSplitViewItem(sidebarItem)
        splitController.addSplitViewItem(contentItem)
        splitController.splitView.autosaveName = "swift.native-notes"

        window = NSWindow(
            contentRect: NSRect(x: 0, y: 0, width: 980, height: 680),
            styleMask: [.titled, .closable, .miniaturizable, .resizable, .fullSizeContentView],
            backing: .buffered,
            defer: false
        )
        window.delegate = self
        window.isReleasedWhenClosed = false
        window.contentMinSize = NSSize(width: 640, height: 420)
        window.backgroundColor = .windowBackgroundColor
        window.isOpaque = true
        window.toolbarStyle = .unified
        window.titlebarSeparatorStyle = .none
        window.contentViewController = splitController
        let toolbar = NSToolbar(identifier: "swift.native-notes.toolbar")
        toolbar.delegate = self
        toolbar.displayMode = .iconOnly
        toolbar.allowsUserCustomization = false
        toolbar.autosavesConfiguration = false
        window.toolbar = toolbar
        // Assigning the split controller can cause AppKit to adopt the
        // controller's fitting size. Restore the same initial content size as
        // the Wails NativeWindow after the complete hierarchy is attached.
        window.setContentSize(NSSize(width: 980, height: 680))
        window.center()

        outlineView.reloadData()
        outlineView.expandItem(section)
        if let first = allFiles.first {
            let row = outlineView.row(forItem: first)
            if row >= 0 {
                outlineView.selectRowIndexes(IndexSet(integer: row), byExtendingSelection: false)
            }
            try open(first)
        }
    }

    private func buildMenuBar() {
        let mainMenu = NSMenu()
        let appItem = NSMenuItem()
        mainMenu.addItem(appItem)
        let appMenu = NSMenu()
        appMenu.addItem(withTitle: "Quit Native Notes Swift", action: #selector(NSApplication.terminate(_:)), keyEquivalent: "q")
        appItem.submenu = appMenu

        let fileItem = NSMenuItem()
        mainMenu.addItem(fileItem)
        let fileMenu = NSMenu(title: "File")
        let saveItem = fileMenu.addItem(withTitle: "Save", action: #selector(saveAction(_:)), keyEquivalent: "s")
        saveItem.target = self
        fileMenu.addItem(.separator())
        let closeItem = fileMenu.addItem(withTitle: "Close", action: #selector(NSWindow.performClose(_:)), keyEquivalent: "w")
        closeItem.target = nil
        fileItem.submenu = fileMenu

        let editItem = NSMenuItem()
        mainMenu.addItem(editItem)
        let editMenu = NSMenu(title: "Edit")
        for (title, action, key) in [
            ("Undo", #selector(UndoManager.undo), "z"),
            ("Redo", #selector(UndoManager.redo), "Z"),
            ("Cut", #selector(NSText.cut(_:)), "x"),
            ("Copy", #selector(NSText.copy(_:)), "c"),
            ("Paste", #selector(NSText.paste(_:)), "v"),
            ("Select All", #selector(NSText.selectAll(_:)), "a"),
        ] {
            editMenu.addItem(withTitle: title, action: action, keyEquivalent: key)
        }
        editItem.submenu = editMenu

        let viewItem = NSMenuItem()
        mainMenu.addItem(viewItem)
        let viewMenu = NSMenu(title: "View")
        let toggle = viewMenu.addItem(withTitle: "Toggle Sidebar", action: #selector(NSSplitViewController.toggleSidebar(_:)), keyEquivalent: "s")
        toggle.keyEquivalentModifierMask = [.control, .command]
        toggle.target = splitController
        viewItem.submenu = viewMenu
        NSApp.mainMenu = mainMenu
    }

    private func buildStatusItem() {
        statusItem = NSStatusBar.system.statusItem(withLength: NSStatusItem.squareLength)
        statusItem.button?.image = NSImage(systemSymbolName: "note.text", accessibilityDescription: "Native Notes")
        let menu = NSMenu()
        let openItem = menu.addItem(withTitle: "Open Native Notes", action: #selector(showWindow(_:)), keyEquivalent: "")
        openItem.target = self
        let saveItem = menu.addItem(withTitle: "Save", action: #selector(saveAction(_:)), keyEquivalent: "")
        saveItem.target = self
        menu.addItem(.separator())
        menu.addItem(withTitle: "Quit", action: #selector(NSApplication.terminate(_:)), keyEquivalent: "")
        statusItem.menu = menu
    }

    @objc private func showWindow(_ sender: Any?) {
        window.makeKeyAndOrderFront(nil)
        NSApp.activate(ignoringOtherApps: true)
        window.makeFirstResponder(textView)
    }

    func windowShouldClose(_ sender: NSWindow) -> Bool {
        sender.orderOut(nil)
        return false
    }

    @objc private func saveAction(_ sender: Any?) {
        do { try save() }
        catch { present(error) }
    }

    private func open(_ file: NoteFile) throws {
        if dirty && currentFile !== file { try save() }
        isLoading = true
        textView.string = try String(contentsOf: file.url, encoding: .utf8)
        currentFile = file
        dirty = false
        isLoading = false
        saveToolbarItem?.isEnabled = false
        saveToolbarItem?.badge = nil
        window.title = "\(file.name) — Native Notes"
        if !benchmark.startsHidden { window.makeFirstResponder(textView) }
    }

    private func save() throws {
        guard let file = currentFile else { return }
        try textView.string.write(to: file.url, atomically: true, encoding: .utf8)
        dirty = false
        saveToolbarItem?.isEnabled = false
        saveToolbarItem?.badge = nil
    }

    func textDidChange(_ notification: Notification) {
        guard !isLoading else { return }
        dirty = true
        saveToolbarItem?.isEnabled = true
    }

    func controlTextDidChange(_ obj: Notification) {
        guard let search = obj.object as? NSSearchField else { return }
        let query = search.stringValue.trimmingCharacters(in: .whitespacesAndNewlines).lowercased()
        section.files = query.isEmpty ? allFiles : allFiles.filter { $0.name.lowercased().contains(query) }
        outlineView.reloadData()
        outlineView.expandItem(section)
    }

    func outlineView(_ outlineView: NSOutlineView, numberOfChildrenOfItem item: Any?) -> Int {
        item == nil ? 1 : (item as? SidebarSection)?.files.count ?? 0
    }

    func outlineView(_ outlineView: NSOutlineView, child index: Int, ofItem item: Any?) -> Any {
        if item == nil { return section }
        return (item as! SidebarSection).files[index]
    }

    func outlineView(_ outlineView: NSOutlineView, isItemExpandable item: Any) -> Bool {
        item is SidebarSection
    }

    func outlineView(_ outlineView: NSOutlineView, isGroupItem item: Any) -> Bool {
        item is SidebarSection
    }

    func outlineView(_ outlineView: NSOutlineView, viewFor tableColumn: NSTableColumn?, item: Any) -> NSView? {
        let identifier = NSUserInterfaceItemIdentifier(item is SidebarSection ? "section" : "file")
        let cell = (outlineView.makeView(withIdentifier: identifier, owner: self) as? NSTableCellView) ?? {
            let view = NSTableCellView()
            view.identifier = identifier
            let image = NSImageView()
            image.translatesAutoresizingMaskIntoConstraints = false
            let label = NSTextField(labelWithString: "")
            label.translatesAutoresizingMaskIntoConstraints = false
            label.lineBreakMode = .byTruncatingTail
            view.imageView = image
            view.textField = label
            view.addSubview(image)
            view.addSubview(label)
            NSLayoutConstraint.activate([
                image.leadingAnchor.constraint(equalTo: view.leadingAnchor, constant: 2),
                image.centerYAnchor.constraint(equalTo: view.centerYAnchor),
                image.widthAnchor.constraint(equalToConstant: 16),
                image.heightAnchor.constraint(equalToConstant: 16),
                label.leadingAnchor.constraint(equalTo: image.trailingAnchor, constant: 7),
                label.trailingAnchor.constraint(equalTo: view.trailingAnchor, constant: -4),
                label.centerYAnchor.constraint(equalTo: view.centerYAnchor),
            ])
            return view
        }()
        if let group = item as? SidebarSection {
            cell.imageView?.image = nil
            cell.textField?.stringValue = group.title.uppercased()
            cell.textField?.font = NSFont.systemFont(ofSize: 11, weight: .semibold)
            cell.textField?.textColor = .secondaryLabelColor
        } else if let file = item as? NoteFile {
            cell.imageView?.image = NSImage(systemSymbolName: "doc.plaintext", accessibilityDescription: nil)
            cell.textField?.stringValue = file.name.replacingOccurrences(of: ".txt", with: "")
            cell.textField?.font = NSFont.systemFont(ofSize: 13)
            cell.textField?.textColor = .labelColor
        }
        return cell
    }

    func outlineViewSelectionDidChange(_ notification: Notification) {
        let row = outlineView.selectedRow
        guard row >= 0, let file = outlineView.item(atRow: row) as? NoteFile else { return }
        do { try open(file) }
        catch { present(error) }
    }

    func toolbarAllowedItemIdentifiers(_ toolbar: NSToolbar) -> [NSToolbarItem.Identifier] {
        [.toggleSidebar, trackingSeparatorIdentifier, searchIdentifier, .flexibleSpace, saveIdentifier]
    }

    func toolbarDefaultItemIdentifiers(_ toolbar: NSToolbar) -> [NSToolbarItem.Identifier] {
        [.toggleSidebar, trackingSeparatorIdentifier, searchIdentifier, .flexibleSpace, saveIdentifier]
    }

    func toolbar(_ toolbar: NSToolbar, itemForItemIdentifier identifier: NSToolbarItem.Identifier, willBeInsertedIntoToolbar flag: Bool) -> NSToolbarItem? {
        switch identifier {
        case .toggleSidebar:
            let item = NSToolbarItem(itemIdentifier: identifier)
            item.label = "Sidebar"
            item.toolTip = "Show or hide the sidebar"
            item.image = NSImage(systemSymbolName: "sidebar.left", accessibilityDescription: item.label)
            item.target = splitController
            item.action = #selector(NSSplitViewController.toggleSidebar(_:))
            return item
        case trackingSeparatorIdentifier:
            return NSTrackingSeparatorToolbarItem(
                identifier: identifier,
                splitView: splitController.splitView,
                dividerIndex: 0
            )
        case searchIdentifier:
            let item = NSSearchToolbarItem(itemIdentifier: identifier)
            item.label = "Search Files"
            item.searchField.placeholderString = "Search Files"
            item.searchField.delegate = self
            return item
        case saveIdentifier:
            let item = NSToolbarItem(itemIdentifier: identifier)
            item.label = "Save"
            item.toolTip = "Save the current text file"
            item.image = NSImage(systemSymbolName: "square.and.arrow.down", accessibilityDescription: item.label)
            item.isBordered = true
            item.isEnabled = false
            item.target = self
            item.action = #selector(saveAction(_:))
            saveToolbarItem = item
            return item
        default:
            return nil
        }
    }

    private func present(_ error: Error) {
        let alert = NSAlert(error: error)
        if window.isVisible { alert.beginSheetModal(for: window) }
        else { alert.runModal() }
    }

    private func signalBenchmarkReady() {
        guard let path = benchmark.readyFile else { return }
        do { try Data("ready\n".utf8).write(to: URL(fileURLWithPath: path)) }
        catch { fputs("benchmark ready file: \(error)\n", stderr) }
    }

    private func scheduleBenchmarkExit() {
        guard let milliseconds = benchmark.autoQuitMilliseconds, milliseconds > 0 else { return }
        DispatchQueue.main.asyncAfter(deadline: .now() + .milliseconds(milliseconds)) {
            NSApp.terminate(nil)
        }
    }
}

private let application = NSApplication.shared
private let delegate = AppDelegate()
application.delegate = delegate
application.run()
