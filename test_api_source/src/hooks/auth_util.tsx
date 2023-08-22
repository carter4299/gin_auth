export function isAscii(str: string) {
    return /^[\x00-\x7F]*$/.test(str);
}

export function validateUsername(username: string) {
    const length = username.length;
    const noSpaces = !/\s/.test(username);
    return length >= 8 && length <= 16 && noSpaces && isAscii(username);
}

export function validatePassword(password: string) {
    const length = password.length;
    const noSpaces = !/\s/.test(password);
    const hasSpecialChar = /[ `!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]/.test(password);
    const hasNumber = /\d/.test(password);
    return length >= 8 && length <= 24 && noSpaces && hasSpecialChar && hasNumber && isAscii(password);
}

export const base64Encode = (str: string) => {
    return btoa(encodeURIComponent(str).replace(/%([0-9A-F]{2})/g, (match, p1) => {
        return String.fromCharCode(Number('0x' + p1));
    }));
}


export const base64EncodeToBytes = (str: string) => {
    const encoded = btoa(encodeURIComponent(str).replace(/%([0-9A-F]{2})/g, (match, p1) => {
        return String.fromCharCode(Number('0x' + p1));
    }));
    const charList = encoded.split('');
    const uintArray = [];
    for (let i = 0; i < charList.length; i++) {
        uintArray.push(charList[i].charCodeAt(0));
    }
    return new Uint8Array(uintArray);
}