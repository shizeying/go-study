document.getElementById('cmdForm').addEventListener('submit', async (event) => {
  event.preventDefault();

  const level = document.getElementById('level').value;
  const sourceSshHost = document.getElementById('sourceSshHost').value;
  const targetSshHost = document.getElementById('targetSshHost').value;
  const passwd = document.getElementById('passwd').value;
  const dockerRunScript = document.getElementById('dockerRunScript').value;
  const script = document.getElementById('script').value;
  const dockerName = document.getElementById('dockerName').value;
  const dockerImage = document.getElementById('dockerImage').value;
  const enable = document.getElementById('enable').value;
  const compression = document.getElementById('compression').value;
  const enable2 = document.getElementById('enable2').value;

  const response = await fetch('/api/execute', {
    method: 'POST', headers: {
      'Content-Type': 'application/json'
    }, body: JSON.stringify({
      level,
      sourceSshHost,
      targetSshHost,
      passwd,
      dockerRunScript,
      script,
      dockerName,
      dockerImage,
      enable,
      compression,
      enable2
    })
  });

  const output = await response.text();
  document.getElementById('output').innerText = output;
});