import Ember from 'ember';
import config from './config/environment';

var Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.resource('home', { path: '/' });

  this.resource('share', { path: '/delen' });

  this.resource('disclaimer', { path: '/voorwaarden' });

  this.resource('input', { path: '/invullen/:input' });

  this.resource('check', { path: '/controlleren' });

  this.resource('preview', { path: '/voorbeeld' });

  this.resource('complete', { path: '/klaar' });

});


export default Router;
