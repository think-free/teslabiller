/**
 * Copyright (c) Facebook, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

'use strict';
import Colors from './Colors';
import {Node} from 'react';
import {Text, StyleSheet, ImageBackground} from 'react-native';
import React from 'react';

const Header = (): Node => (
  <ImageBackground
    accessibilityRole={'image'}
    source={require('../images/logo.png')}
    style={styles.background}
    imageStyle={styles.logo}>
    <Text style={styles.text}>TESLA Biller</Text>
  </ImageBackground>
);

const styles = StyleSheet.create({
  background: {
    paddingBottom: 20,
    paddingTop: 20,
    paddingHorizontal: 32,
    backgroundColor: Colors.secondary,
    borderBottomColor: Colors.blue,
    borderBottomWidth: 3
  },
  logo: {
    opacity: 0.5,
    overflow: 'visible',
    resizeMode: 'cover',
    marginLeft: -256,
    marginBottom: -256,
  },
  text: {
    fontSize: 40,
    fontWeight: '100',
    textAlign: 'center',
    color: Colors.blue,
    fontFamily: 'tesla'
  },
});

export default Header;
