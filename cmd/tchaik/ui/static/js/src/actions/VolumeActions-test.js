var rewire = require('rewire');

describe('VolumeActions', function() {
  var VolumeActions, VolumeConstants;
  var handleViewActionSpy;

  beforeEach(function() {
    React = require('react/addons');
    TestUtils = React.addons.TestUtils;

    VolumeActions = rewire('./VolumeActions.js');
    VolumeConstants = rewire('../constants/VolumeConstants.js');

    handleViewActionSpy = sinon.spy();
    VolumeActions.__Rewire__(
      "AppDispatcher",
      {handleViewAction: handleViewActionSpy}
    );
  });

  describe('volume function', function() {
    describe('when called', function() {
      var volume;

      beforeEach(function() {
        volume = 12;
        VolumeActions.volume(volume);
      });

      it('should call handleViewAction on the AppDispatcher', function() {
        expect(handleViewActionSpy).to.have.been.calledWith({
          actionType: VolumeConstants.SET_VOLUME,
          volume: volume
        });
      });
    });
  });

  describe('toggleVolumeMute function', function() {
    describe('when called', function() {
      beforeEach(function() {
        VolumeActions.toggleVolumeMute();
      });

      it('should call handleViewAction on the AppDispatcher', function() {
        expect(handleViewActionSpy).to.have.been.calledWith({
          actionType: VolumeConstants.TOGGLE_VOLUME_MUTE,
        });
      });
    });
  });
});
