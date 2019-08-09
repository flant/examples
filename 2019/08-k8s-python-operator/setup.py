from sys import version_info

from setuptools import find_packages, setup

if version_info[:2] < (3, 5):
    raise RuntimeError(
        'Unsupported python version %s.' % '.'.join(version_info)
    )


_NAME = 'copyrator'
setup(
    name=_NAME,
    version='0.0.1',
    packages=find_packages(),
    classifiers=[
        'Development Status :: 3 - Alpha',
        'Programming Language :: Python',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.5',
        'Programming Language :: Python :: 3.6',
        'Programming Language :: Python :: 3.7',
    ],
    author='Flant',
    author_email='maksim.nabokikh@flant.com',
    include_package_data=True,
    install_requires=[
        'kubernetes==9.0.0',
    ],
    entry_points={
        'console_scripts': [
            '{0} = {0}.cli:main'.format(_NAME),
        ]
    }
)
