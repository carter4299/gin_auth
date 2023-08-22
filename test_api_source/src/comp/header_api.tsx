const fetchPing = async () => {
    try {
      const response = await fetch("http://localhost:8080/api/go/other/username");
      if (response.status !== 200) {
        throw new Error("Unauthorized");
      }
  
      const data = await response.json();
      console.log(data);
    } catch (error) {
      console.error("Failed to ping:", error);
    }
  };
  
fetchPing();