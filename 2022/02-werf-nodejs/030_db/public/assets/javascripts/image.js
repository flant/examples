import Logo from '../images/werf-logo.svg';
import '../stylesheets/image.css';
import '../stylesheets/style.css';

window.onload = function () {
  const btn = document.getElementById('show-image-btn');
  btn.addEventListener('click', (_) => {
    fetch(Logo)
      .then((data) => data.text())
      .then((html) => {
        const svgContainer = document.getElementById('container');
        svgContainer.insertAdjacentHTML('beforeend', html);
        const svg = svgContainer.getElementsByTagName('svg')[0];
        svg.setAttribute('id', 'image');
        btn.remove();
      });
  });
};
