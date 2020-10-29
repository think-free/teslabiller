/**
 * Update manager
 *
 * @format
 * @flow strict-local
 */

import React, { Component } from 'react';

import {
   StyleSheet,
   View,
   Text,
   StatusBar,
   SafeAreaView,
   Image,
   PermissionsAndroid,
   TextInput,
   ScrollView,
   TouchableOpacity,
   ToastAndroid,
   BackHandler
} from 'react-native';

import DatePicker from 'react-native-datepicker'

var RNFS = require('react-native-fs');

import Header from './components/Header';
import Colors from './components/Colors';

/* Request permitions */

const requestReadStoragePermission = () => {
   try {
      const granted = PermissionsAndroid.request(
         PermissionsAndroid.PERMISSIONS.READ_EXTERNAL_STORAGE,
         {
            title: "External storage permission",
            message:
               "I need to access to your storage to read files",
         }
      );
      if (granted === PermissionsAndroid.RESULTS.GRANTED) {
         console.log("You can read from the storage");
      } else {
         console.log("Storage permission denied");
      }
   } catch (err) {
      console.warn(err);
   }
};

const requestWriteStoragePermission = () => {

   try {
      const granted = PermissionsAndroid.request(
         PermissionsAndroid.PERMISSIONS.WRITE_EXTERNAL_STORAGE,
         {
            title: "External storage permission",
            message:
               "I need to access to your storage to write files",
         }
      );
      if (granted === PermissionsAndroid.RESULTS.GRANTED) {
         console.log("You can write to the storage");
      } else {
         console.log("Storage permission denied");
      }
   } catch (err) {
      console.warn(err);
   }
};

/* Application */

class App extends Component {
   

   constructor(props) {
      super(props);

      this.state = {
         data: {},
         isLoading: true,

         startCounter: "",
         endCounter: "",
         dateStart: "2020-10-01",
         dateEnd: "2020-11-01",
         e1: "",
         e2: "",
         e3: "",
         t1: "",
         t2: "",
         t3: "",
      };

      this.backButtonHanlder = this.backButtonHanlder.bind(this);
   }

   componentDidMount() {

      BackHandler.addEventListener(
         "hardwareBackPress",
         this.backButtonHanlder
      );

      requestReadStoragePermission();
      requestWriteStoragePermission();
      this.loadConfig();
   } 

   backButtonHanlder(){
      this.setState({ data: {} })
      return true
   }

   /* Configuration */

   loadConfig(){

      var path = RNFS.ExternalStorageDirectoryPath + "/.teslabiller.json"
      RNFS.readFile(path)
      .then((content) => {
         var json = JSON.parse(content);
         this.setState({ startCounter: json.startCounter });
         this.setState({ endCounter: json.endCounter });
         this.setState({ dateStart: json.dateStart });
         this.setState({ dateEnd: json.dateEnd });
         this.setState({ e1: json.e1 });
         this.setState({ e2: json.e2 });
         this.setState({ e3: json.e3 });
         this.setState({ t1: json.t1 });
         this.setState({ t2: json.t2 });
         this.setState({ t3: json.t3 });
      })
      .catch((error) => {
         RNFS.writeFile(path, "{}")
            .catch((error) => console.log(error))
      })
      .finally(() => {
         this.setState({ isLoading: false });
      });
   }

   saveConfig(config){

      var path = RNFS.ExternalStorageDirectoryPath + "/.teslabiller.json"
      
      RNFS.writeFile(path, JSON.stringify(config))
         .catch((error) => console.log("Oups : " + error))
   }  

   /* Request server */

   calculate() {

      this.saveConfig(this.state)

      var url = [];
      url.push(
         "http://172.16.10.110:5000/getBillRequest",
         "?dateStart=",
         this.state.dateStart,
         "&dateEnd=",
         this.state.dateEnd,
         "&startCounter=",
         this.state.startCounter,
         "&endCounter=",
         this.state.endCounter,
         "&tax=21",
         "&fix=2",
         "&pE1=",
         this.state.t1,
         "&pE2=",
         this.state.t2,
         "&pE3=",
         this.state.t3,
         "&cE1=",
         this.state.e1,
         "&cE2=",
         this.state.e2,
         "&cE3=",
         this.state.e3,
      );
      var urlStr =  url.join("");

      console.log(urlStr)

      fetch(urlStr)
         .then((response) => response.json())
         .then((json) => {
            if (json.message === undefined){
               this.setState({ data: json });
               console.log(json)
            } else {
               console.log(json.message)
               ToastAndroid.show(json.message, ToastAndroid.LONG);
            }            
         })
   }

   /* Render */

   render() {

      const { data } = this.state;

      if (Object.keys(data).length === 0) {

         console.log("Rendering input")
         return this.renderInput()

      } else {

         console.log("Rendering output")
         return this.renderOutput()
      }      
   }

   renderInput(){
      
      const { isLoading, startCounter, endCounter, e1, e2, e3, t1, t2, t3 } = this.state;

      return (

         <>
            {isLoading ?

               /* Loading screen */
               <View style={styles.loadingContainer}>
                  <Image source={require('./images/logo.png')} style={styles.loading} />
               </View> : (

                  <>
                     <StatusBar barStyle="dark-content" backgroundColor={Colors.blue} />
                     <SafeAreaView style={styles.main}>

                        <Header />

                        <ScrollView style={styles.scrollView}>

                           {/* Date selection */}

                           <View style={styles.spacer} />
                           <Text style={styles.text}>Date interval</Text>
                           <View style={styles.spacer} />

                           <View style={{ flex: 1, flexDirection: 'row', justifyContent: 'space-around' }}>
                              <DatePicker
                                 style={{ width: "40%" }}
                                 date={this.state.dateStart}
                                 mode="date"
                                 placeholder="Starting date"
                                 format="YYYY-MM-DD"
                                 minDate="2020-10-01"
                                 maxDate="2035-10-01"
                                 confirmBtnText="Confirm"
                                 cancelBtnText="Cancel"
                                 iconSource={require('./images/dateStart.png')}
                                 customStyles={{
                                    dateIcon: {
                                       position: 'absolute',
                                       left: 0,
                                       top: 4,
                                       marginLeft: 0
                                    },
                                    dateInput: {
                                       marginLeft: 36,
                                    },
                                    dateText: {
                                       color: Colors.alt1
                                    }
                                 }}
                                 onDateChange={(date) => { this.setState({ dateStart: date }) }}
                              />

                              <DatePicker
                                 style={{ width: "40%" }}
                                 date={this.state.dateEnd}
                                 mode="date"
                                 placeholder="Ending date"
                                 format="YYYY-MM-DD"
                                 minDate="2020-10-01"
                                 maxDate="2035-10-01"
                                 confirmBtnText="Confirm"
                                 cancelBtnText="Cancel"
                                 iconSource={require('./images/dateEnd.png')}
                                 customStyles={{
                                    dateIcon: {
                                       position: 'absolute',
                                       left: 0,
                                       top: 4,
                                       marginLeft: 0
                                    },
                                    dateInput: {
                                       marginLeft: 36
                                    },
                                    dateText: {
                                       color: Colors.alt1
                                    }
                                 }}
                                 onDateChange={(date) => { this.setState({ dateEnd: date }) }}
                              />
                           </View>

                           <View style={styles.hrLine} />

                           {/* Counter */}

                           <View style={styles.spacer} />
                           <Text style={styles.text}>Counter interval</Text>
                           <View style={styles.spacer} />

                           <View style={{ flex: 1, flexDirection: 'row', justifyContent: 'space-around' }}>
                              <TextInput
                                 style={{ height: 40, width: "40%", borderColor: Colors.secondary, borderWidth: 1, color: Colors.alt1 }}
                                 onChangeText={text => this.setState({ startCounter: text })}
                                 value={startCounter}
                                 keyboardType='numeric'
                                 textAlign='right'
                                 placeholder='kWh'
                                 placeholderTextColor={Colors.alt1}
                              />

                              <TextInput
                                 style={{ height: 40, width: "40%", borderColor: Colors.secondary, borderWidth: 1, color: Colors.alt1 }}
                                 onChangeText={text => this.setState({ endCounter: text })}
                                 value={endCounter}
                                 keyboardType='numeric'
                                 textAlign='right'
                                 placeholder='kWh'
                                 placeholderTextColor={Colors.alt1}
                              />
                           </View>

                           <View style={styles.hrLine} />

                           {/* Energy price */}

                           <View style={styles.spacer} />
                           <Text style={styles.text}>Energy price</Text>
                           <View style={styles.spacer} />

                           <View style={{ flex: 1, flexDirection: 'row', justifyContent: 'space-around' }}>
                              <TextInput
                                 style={{ height: 40, width: "30%", borderColor: Colors.secondary, borderWidth: 1, color: Colors.alt1 }}
                                 onChangeText={text => this.setState({ e1: text })}
                                 value={e1}
                                 keyboardType='numeric'
                                 textAlign='right'
                                 placeholder='P1'
                                 placeholderTextColor={Colors.alt1}
                              />

                              <TextInput
                                 style={{ height: 40, width: "30%", borderColor: Colors.secondary, borderWidth: 1, color: Colors.alt1 }}
                                 onChangeText={text => this.setState({ e2: text })}
                                 value={e2}
                                 keyboardType='numeric'
                                 textAlign='right'
                                 placeholder='P2'
                                 placeholderTextColor={Colors.alt1}
                              />

                              <TextInput
                                 style={{ height: 40, width: "30%", borderColor: Colors.secondary, borderWidth: 1, color: Colors.alt1 }}
                                 onChangeText={text => this.setState({ e3: text })}
                                 value={e3}
                                 keyboardType='numeric'
                                 textAlign='right'
                                 placeholder='P3'
                                 placeholderTextColor={Colors.alt1}
                              />
                           </View>

                           <View style={styles.hrLine} />

                           {/* Toll price */}

                           <View style={styles.spacer} />
                           <Text style={styles.text}>Toll price</Text>
                           <View style={styles.spacer} />

                           <View style={{ flex: 1, flexDirection: 'row', justifyContent: 'space-around' }}>
                              <TextInput
                                 style={{ height: 40, width: "30%", borderColor: Colors.secondary, borderWidth: 1, color: Colors.alt1 }}
                                 onChangeText={text => this.setState({ t1: text })}
                                 value={t1}
                                 keyboardType='numeric'
                                 textAlign='right'
                                 placeholder='P1'
                                 placeholderTextColor={Colors.alt1}
                              />

                              <TextInput
                                 style={{ height: 40, width: "30%", borderColor: Colors.secondary, borderWidth: 1, color: Colors.alt1 }}
                                 onChangeText={text => this.setState({ t2: text })}
                                 value={t2}
                                 keyboardType='numeric'
                                 textAlign='right'
                                 placeholder='P2'
                                 placeholderTextColor={Colors.alt1}
                              />

                              <TextInput
                                 style={{ height: 40, width: "30%", borderColor: Colors.secondary, borderWidth: 1, color: Colors.alt1 }}
                                 onChangeText={text => this.setState({ t3: text })}
                                 value={t3}
                                 keyboardType='numeric'
                                 textAlign='right'
                                 placeholder='P3'
                                 placeholderTextColor={Colors.alt1}
                              />
                           </View>

                           <View style={styles.spacerBig} />

                        </ScrollView>

                        <TouchableOpacity style={styles.touch} onPress={() => this.calculate()}>
                           <Image source={require('./images/check.png')} style={styles.image} />
                        </TouchableOpacity>

                     </SafeAreaView>
                  </>

               )}
         </>
      );
   }

   renderOutput(){

      const { data } = this.state;

      return (

         <>
            <StatusBar barStyle="dark-content" backgroundColor={Colors.blue} />
            <SafeAreaView style={styles.main}>

               <Header />

               <ScrollView style={styles.scrollView}>

                  <View style={styles.hrLine} />
                  <View style={styles.spacer} />

                  <Text style={styles.text}>Total</Text>
                  <View style={styles.spacer} />
                  <Text style={styles.total}>{data.totalCostTax} €</Text>
                  <Text style={styles.totalSub}>{data.totalkWh} kWh</Text>

                  <View style={styles.hrLine} />
                  <View style={styles.spacer} />

                  <Text style={styles.text}>Car</Text>
                  <View style={styles.spacerSmall} />
                  <Text style={styles.totalAlt}>{data.carTotalCostTax} €</Text>
                  <Text style={styles.totalSubAlt}>{data.carTotalkWh} kWh</Text>

                  <View style={styles.spacer} />

                  <View style={{ flex: 1, flexDirection: 'row', justifyContent: 'space-between' }}>
                     <Text style={styles.costDetail}>P1</Text>
                     <Text style={styles.costDetail}>{data.carP1kWh} kWh</Text>
                     <Text style={styles.costDetail}>{data.carP1Cost} €</Text>
                  </View>
                  <View style={{ flex: 1, flexDirection: 'row', justifyContent: 'space-between' }}>
                     <Text style={styles.costDetail}>P2</Text>
                     <Text style={styles.costDetail}>{data.carP2kWh} kWh</Text>
                     <Text style={styles.costDetail}>{data.carP2Cost} €</Text>
                  </View>
                  <View style={{ flex: 1, flexDirection: 'row', justifyContent: 'space-between' }}>
                     <Text style={styles.costDetail}>P3</Text>
                     <Text style={styles.costDetail}>{data.carP3kWh} kWh</Text>
                     <Text style={styles.costDetail}>{data.carP3Cost} €</Text>
                  </View>

                  <View style={styles.hrLine} />
                  <View style={styles.spacer} />

                  <Text style={styles.text}>Other</Text>
                  <View style={styles.spacerSmall} />
                  <Text style={styles.totalAlt}>{data.otherCostTax} €</Text>
                  <Text style={styles.totalSubAlt}>{data.otherTotalkWh} kWh</Text>

               </ScrollView>
            </SafeAreaView>
         </>
      );
   }
};

const styles = StyleSheet.create({

   spacerSmall: {
      height: 10
   },

   spacer:{
      height: 20
   },

   spacerBig: {
      height: 150
   },

   hrLine: {
      marginTop: 30,
      marginLeft:10,
      marginRight: 10,
      backgroundColor: Colors.blue,
      height:1,
      opacity: 0.5
   },

   loadingContainer: {

      backgroundColor: Colors.main,
      height: "100%",
      width: "100%",
      justifyContent: 'center',
      alignItems: 'center',
   },
   loading: {
      justifyContent: 'center',
      alignItems: 'center',
      height: 200,
      width: 200,
   },
   text: {
      fontSize: 20,
      fontWeight: '400',
      textAlign: 'center',
      color: Colors.secondary,
   },

   main: {
      backgroundColor: Colors.main,
      height: "100%",
   },
   scrollView: {
      backgroundColor: Colors.main,
      padding: 10,
   },

   image: {
      position: 'absolute',
      top: 0,
      right: 0,
      height: 60,
      width: 60,
   },
   touch: {
      position: 'absolute',
      bottom: 30,
      right: 30,
      height: 60,
      width: 60,
      backgroundColor: Colors.secondary,
      borderRadius: 60
   },

   total: {
      fontSize: 40,
      fontWeight: '400',
      textAlign: 'center',
      color: Colors.blue,
   },
   totalSub: {
      fontSize: 20,
      fontWeight: '400',
      textAlign: 'center',
      color: Colors.alt1,
      marginTop: 10,
   },

   totalAlt: {
      fontSize: 20,
      fontWeight: '400',
      textAlign: 'center',
      color: Colors.blue,
   },
   totalSubAlt: {
      fontSize: 15,
      fontWeight: '400',
      textAlign: 'center',
      color: Colors.alt1,
      marginTop: 5,
   },

   costDetail:{
      fontSize: 16,
      color: Colors.alt1,
      textAlign: 'right'
   }
});

export default App;
