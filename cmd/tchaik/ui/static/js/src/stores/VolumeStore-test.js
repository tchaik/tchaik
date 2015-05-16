/* jshint expr:true */

var rewire = require('rewire');

describe('VolumeStore', function() {
  var VolumeConstants = require('../constants/VolumeConstants.js');
  var VolumeStore;
  var registerSpy;
  var callback;


  var actionToggleMute = {
    source: 'VIEW_ACTION',
    action: {
      actionType: VolumeConstants.TOGGLE_VOLUME_MUTE
    }
  };

  beforeEach(function() {
    // we have to clear localstorage, otherwise we have leaks from other tests
    localStorage.clear();

    // We can't quite get the rewriting working, due to the fact that
    // AppDispatcher is called before anything is exported from the stores. The
    // only way we can mock the register function is by applying the spy on the
    // raw module using require rather than rewire.
    var AppDispatcher = require('../dispatcher/AppDispatcher.js');
    registerSpy = sinon.spy(function(cb) {
      callback = cb;
    });
    AppDispatcher.register = registerSpy;
    VolumeStore = rewire('./VolumeStore.js');
  });

  describe('on initialisation', function() {
    it('registers a callback with the dispatcher', function() {
      expect(registerSpy).to.have.been.called;
    });

    describe('the volume', function() {
      it('should be at it\'s default value of 0.75', function() {
        expect(VolumeStore.getVolume()).to.equal(0.75);
      });
    });

    describe('the mute state', function() {
      it('should be unmuted', function() {
        expect(VolumeStore.getVolumeMute()).not.to.be.ok;
      });
    });
  });

  describe('when triggering the SET_VOLUME action', function() {
    var volume = 1.0;
    beforeEach(function() {
      callback({
        source: 'VIEW_ACTION',
        action: {
          actionType: VolumeConstants.SET_VOLUME,
          volume: volume
        }
      });
    });

    describe('the volume', function() {
      it('should be the value that was set in the action', function() {
        expect(VolumeStore.getVolume()).to.equal(volume);
      });
    });

    describe('the mute status', function() {
      it('should be unmuted', function() {
        expect(VolumeStore.getVolumeMute()).not.to.be.ok;
      });
    });
  });

  describe('triggering the TOGGLE_VOLUME_MUTE action', function() {
    describe('when it is not currently muted', function() {
      beforeEach(function() {
        callback(actionToggleMute);
      });

      describe('the mute status', function() {
        it('should be muted', function() {
          expect(VolumeStore.getVolumeMute()).to.be.ok;
        });
      });

      describe('the volume', function() {
        it('should be 0', function() {
          expect(VolumeStore.getVolume()).to.equal(0);
        });
      });
    });

    describe('when it is currently muted', function() {
      beforeEach(function() {
        callback(actionToggleMute);
        callback(actionToggleMute);
      });

      describe('the mute status', function() {
        it('should be muted', function() {
          expect(VolumeStore.getVolumeMute()).not.to.be.ok;
        });
      });

      describe('the volume', function() {
        it('should be the original value', function() {
          expect(VolumeStore.getVolume()).to.equal(0.75);
        });
      });
    });
  });
});

