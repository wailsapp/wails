export function log(message) {
    // eslint-disable-next-line
    console.log(
        '%c wails dev %c ' + message + ' ',
        'background: #aa0000; color: #fff; border-radius: 3px 0px 0px 3px; padding: 1px; font-size: 0.7rem',
        'background: #009900; color: #fff; border-radius: 0px 3px 3px 0px; padding: 1px; font-size: 0.7rem'
    );
}