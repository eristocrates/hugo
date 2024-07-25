import * as vscode from 'vscode';

export function activate(context: vscode.ExtensionContext) {
    const linkProvider: vscode.DocumentLinkProvider = {
        provideDocumentLinks(document: vscode.TextDocument, token: vscode.CancellationToken) {
            const links: vscode.DocumentLink[] = [];
            const regex = /"caller":"([^"]+):(\d+)"/g;
            const text = document.getText();

            let match;
            while ((match = regex.exec(text)) !== null) {
                const filePath = match[1].replace(/\\/g, '/'); // Replace backslashes with forward slashes
                const lineNumber = match[2];
                const index = match.index + 9; // +9 to account for the "caller":" prefix
                const length = match[0].length - 1; // -1 to exclude the trailing double quote

                const range = new vscode.Range(
                    document.positionAt(index),
                    document.positionAt(index + length)
                );

                const target = vscode.Uri.parse(`command:extension.openFileAtLine?${encodeURIComponent(JSON.stringify({ filePath, lineNumber }))}`);
                links.push(new vscode.DocumentLink(range, target));
            }

            return links;
        }
    };

    context.subscriptions.push(
        vscode.languages.registerDocumentLinkProvider({ scheme: 'file', language: 'json' }, linkProvider)
    );

    context.subscriptions.push(
        vscode.commands.registerCommand('extension.openFileAtLine', (file) => {
            const { filePath, lineNumber } = file;
            const uri = vscode.Uri.file(filePath);
            vscode.workspace.openTextDocument(uri).then(doc => {
                vscode.window.showTextDocument(doc).then(editor => {
                    const line = parseInt(lineNumber) - 1;
                    const position = new vscode.Position(line, 0);
                    editor.selection = new vscode.Selection(position, position);
                    editor.revealRange(new vscode.Range(position, position));
                });
            });
        })
    );
}

export function deactivate() {}
