# Meta Knight

This project explores the use of GraphQL for storing and querying build job metadata.


    curl -g 'http://localhost:12345/graphql?query={builds(project:"1"){project,moduleid}}'