-- go
DELETE FROM deployment_techstack WHERE name = 'Go' AND version = '1.20.14';
DELETE FROM deployment_techstack WHERE name = 'Go' AND version = '1.21.13';
DELETE FROM deployment_techstack WHERE name = 'Go' AND version = '1.22.12';
DELETE FROM deployment_techstack WHERE name = 'Go' AND version = '1.23.12';
DELETE FROM deployment_techstack WHERE name = 'Go' AND version = '1.25.9';
DELETE FROM deployment_techstack WHERE name = 'Go' AND version = '1.26.2';

-- nodejs
DELETE FROM deployment_techstack WHERE name = 'Node.js' AND version = 'v20 LTS';
DELETE FROM deployment_techstack WHERE name = 'Node.js' AND version = 'v21';
DELETE FROM deployment_techstack WHERE name = 'Node.js' AND version = 'v22 LTS';
DELETE FROM deployment_techstack WHERE name = 'Node.js' AND version = 'v23';
DELETE FROM deployment_techstack WHERE name = 'Node.js' AND version = 'v24 LTS';
DELETE FROM deployment_techstack WHERE name = 'Node.js' AND version = 'v25';
DELETE FROM deployment_techstack WHERE name = 'Node.js' AND version = 'v26 LTS';

-- python
DELETE FROM deployment_techstack WHERE name = 'Python' AND version = '3.9.x';
DELETE FROM deployment_techstack WHERE name = 'Python' AND version = '3.10.x';
DELETE FROM deployment_techstack WHERE name = 'Python' AND version = '3.11.x';
DELETE FROM deployment_techstack WHERE name = 'Python' AND version = '3.12.x';
DELETE FROM deployment_techstack WHERE name = 'Python' AND version = '3.13.x';
DELETE FROM deployment_techstack WHERE name = 'Python' AND version = '3.14.x';
