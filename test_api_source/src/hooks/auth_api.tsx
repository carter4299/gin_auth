export async function authenticateUser(encodedUserName: string, encodedUserPass: Uint8Array): Promise<Response> {
    let formData = new FormData();
    formData.append("user_id", encodedUserName);
    formData.append("password", new Blob([encodedUserPass]));
    console.log(formData);

    return await fetch('http://localhost:8080/api/go/auth/login', {
        method: 'POST',
        body: formData,
        credentials: 'include'
    });
}

export async function registerUser(encodedUserName: string, encodedUserPass: Uint8Array): Promise<Response> {
    let formData = new FormData();
    formData.append("user_id", encodedUserName);
    formData.append("password", new Blob([encodedUserPass]));
    console.log(formData);
    return await fetch('http://localhost:8080/api/go/auth/signup', {
        method: 'POST',
        body: formData,
        credentials: 'include'
    });
}